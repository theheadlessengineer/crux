package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	dataplugins "github.com/theheadlessengineer/crux/data/plugins"
	"github.com/theheadlessengineer/crux/internal/app/config"
	"github.com/theheadlessengineer/crux/internal/domain/lockfile"
	"github.com/theheadlessengineer/crux/internal/domain/model"
	"github.com/theheadlessengineer/crux/internal/domain/plugin"
	"github.com/theheadlessengineer/crux/internal/domain/prompt"
	infraconfig "github.com/theheadlessengineer/crux/internal/infrastructure/config"
	"github.com/theheadlessengineer/crux/internal/infrastructure/generator"
	infraplugin "github.com/theheadlessengineer/crux/internal/infrastructure/plugin"
	"github.com/theheadlessengineer/crux/internal/presentation/tui"
)

type newCommandFlags struct {
	outputDir  string
	configFile string
	dryRun     bool
	noPrompt   bool
}

func newNewCommand(_ *config.GlobalConfig, cruxVersion string) *cobra.Command {
	flags := &newCommandFlags{}

	cmd := &cobra.Command{
		Use:   "new <service-name>",
		Short: "Generate a new microservice skeleton",
		Long:  "Initiates the service generation flow for a new microservice.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			if err := model.ValidateServiceName(name); err != nil {
				return &ValidationError{Msg: err.Error()}
			}

			var fileCfg *infraconfig.Config
			if flags.configFile != "" {
				loaded, err := infraconfig.Load(flags.configFile)
				if err != nil {
					return fmt.Errorf("load config: %w", err)
				}
				fileCfg = loaded
			}

			if flags.noPrompt {
				if fileCfg == nil {
					return infraconfig.ErrNoConfigForNoPrompt
				}
				if err := infraconfig.ValidateForNoPrompt(fileCfg); err != nil {
					return err
				}
			}

			if flags.dryRun {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "[dry-run] Would generate service: %s\n", name)
				return nil
			}

			var answers map[string]prompt.Answer
			var selectedPlugins []*plugin.Plugin
			if !flags.noPrompt {
				var err error
				answers, selectedPlugins, err = runPrompt(cmd.InOrStdin(), cmd.OutOrStdout(), name, cruxVersion)
				if err != nil {
					if errors.Is(err, tui.ErrAborted) {
						_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Aborted.")
						return nil
					}
					return err
				}
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\nGenerating %s...\n", name)

			outDir := flags.outputDir
			if outDir == "" {
				outDir = name
			}

			if err := os.MkdirAll(outDir, 0o750); err != nil {
				return fmt.Errorf("create output directory %q: %w", outDir, err)
			}

			genCfg := buildGeneratorConfig(name, cruxVersion, fileCfg, answers, selectedPlugins)
			if err := generator.Generate(cmd.Context(), &genCfg, outDir); err != nil {
				return fmt.Errorf("generate skeleton: %w", err)
			}

			skel := buildSkeleton(name, cruxVersion, fileCfg, answers, selectedPlugins)
			if err := lockfile.Write(outDir, skel); err != nil {
				return fmt.Errorf("write lockfiles: %w", err)
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "✔  skeleton generated\n")
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "✔  .skeleton.json written\n")
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "✔  crux.lock written\n")
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\nNext steps:\n  cd %s && make dev\n  cat docs/TODO.md\n", outDir)
			return nil
		},
	}

	cmd.Flags().StringVar(&flags.outputDir, "output-dir", "",
		"Directory to write the generated service (default: ./<service-name>)")
	cmd.Flags().StringVar(&flags.configFile, "config", "", "Path to a pre-filled configuration YAML file")
	cmd.Flags().BoolVar(&flags.dryRun, "dry-run", false, "Print what would be generated without writing files")
	cmd.Flags().BoolVar(&flags.noPrompt, "no-prompt", false, "Run non-interactively using config file or defaults only")

	return cmd
}

// coreQuestions returns the service-level questions asked before plugin selection.
func coreQuestions() []prompt.Question {
	return []prompt.Question{
		{
			ID:     "language",
			Type:   prompt.QuestionTypeSelect,
			Prompt: "Language + Framework",
			Help: "The programming language and framework for this service. " +
				"All Tier 1 standards are generated for every language.",
			Options: []prompt.Option{
				{Label: "Go + Gin", Value: "go"},
				{Label: "Python + FastAPI", Value: "python"},
				{Label: "Java + Spring Boot", Value: "java"},
				{Label: "Node.js + Express", Value: "node"},
			},
			Default: "go",
		},
		{
			ID:      "team",
			Type:    prompt.QuestionTypeText,
			Prompt:  "Team name",
			Default: "platform",
			Help: "The owning team for this service. Used in generated README, runbooks, and alert" +
				" routing rules. Must be lowercase alphanumeric with hyphens (e.g. payments, data-platform).",
			Validation: prompt.ValidationRule{
				Required: true,
				Pattern:  `^[a-z][a-z0-9-]{0,62}$`,
			},
		},
		{
			ID:     "module",
			Type:   prompt.QuestionTypeText,
			Prompt: "Go module path (e.g. github.com/org/service-name)",
			Help: "The Go module path declared in go.mod. This must be unique across the organisation" +
				" and follow your VCS layout (e.g. github.com/acme/payment-service)." +
				" Leave blank to set it manually later.",
		},
		{
			ID:      "slo_availability",
			Type:    prompt.QuestionTypeText,
			Prompt:  "Availability SLO target (e.g. 99.9)",
			Default: "99.9",
			Help: "The percentage of time this service must be available (e.g. 99.9 = three nines)." +
				" Written into the generated SLO config and used to derive error budget alerts." +
				" Optional — defaults to 99.9%.",
		},
		{
			ID:      "slo_p99_latency_ms",
			Type:    prompt.QuestionTypeNumber,
			Prompt:  "p99 latency target (ms)",
			Default: "500",
			Help: "The 99th-percentile response time budget in milliseconds. Used to generate" +
				" Prometheus alerting rules and SLO dashboards." +
				" A value of 500 means 99% of requests must complete within 500ms.",
			Validation: prompt.ValidationRule{
				Min: 1,
				Max: 60000,
			},
		},
	}
}

// loadAvailablePlugins loads all plugins from the embedded data/plugins/ FS.
func loadAvailablePlugins(cruxVersion string) ([]*plugin.Plugin, error) {
	return infraplugin.LoadFromFS(dataplugins.FS, cruxVersion)
}

// pluginQuestionToPrompt converts a plugin QuestionSpec to a prompt.Question.
// language is the value of the "language" answer (e.g. "go", "python", "java", "node").
// If the spec declares options_by_language, those override the generic options for the
// selected language. Same for default_by_language.
func pluginQuestionToPrompt(qs *plugin.QuestionSpec, pluginName, language string) prompt.Question {
	q := prompt.Question{
		ID:      pluginName + "." + qs.ID,
		Prompt:  qs.Prompt,
		Help:    qs.Help,
		Default: qs.Default,
	}

	// Resolve language-specific default.
	if language != "" && len(qs.DefaultByLang) > 0 {
		if d, ok := qs.DefaultByLang[language]; ok {
			q.Default = d
		}
	}

	// Resolve language-specific options.
	rawOptions := qs.Options
	if language != "" && len(qs.OptionsByLang) > 0 {
		if langOpts, ok := qs.OptionsByLang[language]; ok {
			rawOptions = langOpts
		}
	}

	switch qs.Type {
	case "confirm":
		q.Type = prompt.QuestionTypeConfirm
	case "select":
		q.Type = prompt.QuestionTypeSelect
		q.Options = make([]prompt.Option, len(rawOptions))
		for i, o := range rawOptions {
			q.Options[i] = prompt.Option{Label: o, Value: o}
		}
	case "number":
		q.Type = prompt.QuestionTypeNumber
	default:
		q.Type = prompt.QuestionTypeText
	}
	return q
}

// runPrompt drives the full interactive session in a single TUI pass.
// All plugin questions are included upfront, gated by DependsOn on the _plugins answer.
func runPrompt(
	in io.Reader, out io.Writer, serviceName, cruxVersion string,
) (map[string]prompt.Answer, []*plugin.Plugin, error) {
	allPlugins, err := loadAvailablePlugins(cruxVersion)
	if err != nil {
		allPlugins = nil
	}

	questions := coreQuestions()

	// Determine language from a pre-scan of the config file or default to "go".
	// At prompt time the language answer isn't known yet, so we use a sentinel
	// and resolve per-question options dynamically via a two-pass approach:
	// the language question is always first in coreQuestions(), so we build the
	// full question list with a placeholder language and re-resolve after the
	// session collects the language answer.
	//
	// Simpler approach: build questions with all language variants embedded as
	// separate conditional questions, one per language per option set.
	// We use the DependsOn mechanism: for each plugin question that has
	// options_by_language, we emit one question per language variant, each
	// gated on language == <lang>. Only the matching one will be visible.

	pluginMap := make(map[string]*plugin.Plugin, len(allPlugins))
	if len(allPlugins) > 0 {
		selQ := prompt.Question{
			ID:     "_plugins",
			Type:   prompt.QuestionTypeMultiSelect,
			Prompt: "Select integrations to include",
			Help: "Choose the integrations your service needs. Each selected plugin adds" +
				" pre-configured boilerplate — connection pooling, health checks, metrics," +
				" and infrastructure code. You can add more plugins later with `crux plugin install`.",
			Options: make([]prompt.Option, len(allPlugins)),
		}
		for i, p := range allPlugins {
			label := p.Manifest.Metadata.Name
			if p.Manifest.Metadata.Description != "" {
				label += " — " + p.Manifest.Metadata.Description
			}
			selQ.Options[i] = prompt.Option{Label: label, Value: p.Manifest.Metadata.Name}
			pluginMap[p.Manifest.Metadata.Name] = p
		}
		questions = append(questions, selQ)

		// For each plugin question, emit language-aware variants.
		for _, p := range allPlugins {
			for i := range p.Manifest.Spec.Questions {
				qs := &p.Manifest.Spec.Questions[i]
				pluginGate := prompt.DependsOn{
					And: []prompt.Condition{{QuestionID: "_plugins", Value: p.Manifest.Metadata.Name}},
				}

				if len(qs.OptionsByLang) > 0 {
					// Emit one question per language variant, each gated on language == <lang>.
					for _, lang := range []string{"go", "python", "java", "node"} {
						q := pluginQuestionToPrompt(qs, p.Manifest.Metadata.Name, lang)
						// Suffix the ID so each variant is unique in the graph.
						q.ID = q.ID + "." + lang
						q.DependsOn = &prompt.DependsOn{
							And: append(pluginGate.And, prompt.Condition{QuestionID: "language", Value: lang}),
						}
						questions = append(questions, q)
					}
				} else {
					// No language-specific options — single question gated on plugin selection only.
					q := pluginQuestionToPrompt(qs, p.Manifest.Metadata.Name, "")
					q.DependsOn = &pluginGate
					questions = append(questions, q)
				}
			}
		}
	}

	scanner := bufio.NewScanner(in)
	answers, err := runSession(in, out, scanner, questions, serviceName)
	if err != nil {
		return nil, nil, err
	}

	var selectedPlugins []*plugin.Plugin
	if a, ok := answers["_plugins"]; ok {
		if chosen, ok := a.Value.([]string); ok {
			for _, name := range chosen {
				if p, ok := pluginMap[name]; ok {
					selectedPlugins = append(selectedPlugins, p)
				}
			}
		}
	}

	return answers, selectedPlugins, nil
}

// runSession drives a prompt.Session for the given questions and returns collected answers.
// Uses the Bubbletea TUI when in is a TTY; falls back to plain text otherwise.
func runSession(
	in io.Reader,
	out io.Writer,
	scanner *bufio.Scanner,
	questions []prompt.Question,
	serviceName string,
) (map[string]prompt.Answer, error) {
	graph, err := prompt.NewDecisionGraph(questions, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("build question graph: %w", err)
	}
	session := prompt.NewSession(graph)

	if isReaderTTY(in) && isWriterTTY(out) {
		return tui.Run(session, serviceName)
	}

	// Plain-text fallback (piped input / tests).
	for {
		q := session.NextQuestion()
		if q == nil {
			break
		}

		def := ""
		if q.Default != "" {
			def = fmt.Sprintf(" [%s]", q.Default)
		}
		if len(q.Options) > 0 {
			opts := make([]string, len(q.Options))
			for i, o := range q.Options {
				opts[i] = o.Value
			}
			_, _ = fmt.Fprintf(out, "? %s%s\n  (%s): ", q.Prompt, def, strings.Join(opts, ", "))
		} else {
			_, _ = fmt.Fprintf(out, "? %s%s: ", q.Prompt, def)
		}

		more := scanner.Scan()
		raw := strings.TrimSpace(scanner.Text())

		// EOF — accept default and move on.
		if !more && raw == "" {
			raw = q.Default
		}

		if raw == "b" {
			if backErr := session.Back(); backErr != nil {
				_, _ = fmt.Fprintf(out, "  ! %s\n", backErr)
			}
			continue
		}

		answer, valErr := prompt.Validate(q, raw)
		if valErr != nil {
			if !more {
				// No more input and validation failed — skip this question.
				break
			}
			_, _ = fmt.Fprintf(out, "  ✘ %s\n", valErr)
			continue
		}

		session.Record(q, answer)
	}

	return session.Answers(), nil
}

// isReaderTTY reports whether r is an *os.File connected to an interactive terminal.
func isReaderTTY(r io.Reader) bool {
	f, ok := r.(*os.File)
	if !ok {
		return false
	}
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

// isWriterTTY reports whether w is an *os.File connected to an interactive terminal.
func isWriterTTY(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func buildGeneratorConfig(
	name, cruxVersion string,
	fileCfg *infraconfig.Config,
	answers map[string]prompt.Answer,
	selectedPlugins []*plugin.Plugin,
) generator.Config {
	cfg := generator.Config{
		ServiceName: name,
		Language:    "go",
		Framework:   "gin",
		CLIVersion:  cruxVersion,
		GeneratedAt: time.Now().UTC(),
	}
	if fileCfg != nil {
		if fileCfg.Service.Language != "" {
			cfg.Language = fileCfg.Service.Language
		}
		if fileCfg.Service.Framework != "" {
			cfg.Framework = fileCfg.Service.Framework
		}
		if fileCfg.Service.Team != "" {
			cfg.Team = fileCfg.Service.Team
		}
	}
	if answers != nil {
		if a, ok := answers["language"]; ok {
			if lang, _ := a.Value.(string); lang != "" {
				cfg.Language = lang
				cfg.Framework = defaultFramework(lang)
			}
		}
		if a, ok := answers["team"]; ok {
			cfg.Team, _ = a.Value.(string)
		}
		if a, ok := answers["module"]; ok {
			if m, _ := a.Value.(string); m != "" {
				cfg.Module = m
			}
		}
		// Populate flat answers map for template data.
		cfg.Answers = make(map[string]any, len(answers))
		for k, a := range answers {
			cfg.Answers[k] = a.Value
		}
	}

	// Resolve plugin templates for the selected language.
	for _, p := range selectedPlugins {
		cfg.Plugins = append(cfg.Plugins, generator.SelectedPlugin{
			Name:      p.Manifest.Metadata.Name,
			Templates: p.Manifest.Spec.TemplatesForLang(cfg.Language),
		})
	}

	return cfg
}

// defaultFramework returns the canonical framework name for a given language selection.
func defaultFramework(language string) string {
	switch language {
	case "python":
		return "fastapi"
	case "java":
		return "spring"
	case "node":
		return "express"
	default:
		return "gin"
	}
}

func buildSkeleton(
	name, cruxVersion string,
	fileCfg *infraconfig.Config,
	answers map[string]prompt.Answer,
	selectedPlugins []*plugin.Plugin,
) *lockfile.Skeleton {
	skel := &lockfile.Skeleton{
		CruxVersion: cruxVersion,
		GeneratedAt: time.Now().UTC(),
		Service: lockfile.SkeletonService{
			Name: name,
		},
		Answers:    make(map[string]any),
		Plugins:    []lockfile.PluginEntry{},
		Deviations: []string{},
		Tier1Standards: lockfile.Tier1Standards{
			Enforced:          true,
			DisabledStandards: []string{},
		},
	}

	if fileCfg != nil {
		if fileCfg.Service.Language != "" {
			skel.Service.Language = fileCfg.Service.Language
		}
		if fileCfg.Service.Framework != "" {
			skel.Service.Framework = fileCfg.Service.Framework
		}
		for k, v := range fileCfg.Answers {
			skel.Answers[k] = v
		}
	}

	for k, a := range answers {
		skel.Answers[k] = a.Value
	}

	for _, p := range selectedPlugins {
		skel.Plugins = append(skel.Plugins, lockfile.PluginEntry{
			Name:    p.Manifest.Metadata.Name,
			Version: p.Manifest.Metadata.Version,
		})
	}

	return skel
}
