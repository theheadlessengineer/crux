// Package plugin implements plugin manifest parsing, validation, and loading.
package plugin

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"

	domain "github.com/theheadlessengineer/crux/internal/domain/plugin"
)

const (
	manifestFile = "plugin.yaml"
	pluginsDir   = ".crux/plugins"
)

var semverRE = regexp.MustCompile(`^\d+\.\d+\.\d+$`)

// loader implements domain.Loader.
type loader struct {
	embeddedDirs []string // absolute paths to embedded plugin dirs (injected for testing)
}

// New returns a Loader that searches embeddedDirs and ~/.crux/plugins/.
// Pass nil embeddedDirs to skip the embedded search path.
func New(embeddedDirs []string) domain.Loader {
	return &loader{embeddedDirs: embeddedDirs}
}

// Load discovers, parses, validates, and returns all plugins.
func (l *loader) Load(cruxVersion string) ([]*domain.Plugin, error) {
	dirs, err := l.discoverDirs()
	if err != nil {
		return nil, err
	}

	var plugins []*domain.Plugin
	for _, dir := range dirs {
		p, err := loadOne(dir, cruxVersion)
		if err != nil {
			return nil, fmt.Errorf("plugin %s: %w", dir, err)
		}
		plugins = append(plugins, p)
	}
	return plugins, nil
}

// discoverDirs returns all directories that contain a plugin.yaml.
func (l *loader) discoverDirs() ([]string, error) {
	var dirs []string

	for _, base := range l.embeddedDirs {
		found, err := findPluginDirs(base)
		if err != nil {
			return nil, err
		}
		dirs = append(dirs, found...)
	}

	home, err := os.UserHomeDir()
	if err == nil {
		found, _ := findPluginDirs(filepath.Join(home, pluginsDir))
		dirs = append(dirs, found...)
	}

	return dirs, nil
}

// findPluginDirs walks root and returns every directory containing plugin.yaml.
func findPluginDirs(root string) ([]string, error) {
	var dirs []string
	entries, err := os.ReadDir(root)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		candidate := filepath.Join(root, e.Name())
		if _, err := os.Stat(filepath.Join(candidate, manifestFile)); err == nil {
			dirs = append(dirs, candidate)
		}
	}
	return dirs, nil
}

// loadOne parses and validates the plugin.yaml in dir, then returns a Plugin.
func loadOne(dir, cruxVersion string) (*domain.Plugin, error) {
	m, err := ParseManifest(filepath.Join(dir, manifestFile))
	if err != nil {
		return nil, err
	}
	if err := ValidateManifest(m); err != nil {
		return nil, err
	}
	if err := checkCompatibility(m.Metadata.CruxVersionConstraint, cruxVersion); err != nil {
		return nil, err
	}
	return &domain.Plugin{Manifest: m}, nil
}

// ParseManifest reads and unmarshals a plugin.yaml file.
func ParseManifest(path string) (*domain.Manifest, error) {
	data, err := os.ReadFile(path) //nolint:gosec // path is application-controlled
	if err != nil {
		return nil, fmt.Errorf("read manifest: %w", err)
	}
	var m domain.Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parse manifest: %w", err)
	}
	return &m, nil
}

// ValidateManifest checks required fields, version format, and trust tier.
func ValidateManifest(m *domain.Manifest) error {
	if m.APIVersion != "crux/v1" {
		return fmt.Errorf("apiVersion must be \"crux/v1\", got %q", m.APIVersion)
	}
	if m.Kind != "Plugin" {
		return fmt.Errorf("kind must be \"Plugin\", got %q", m.Kind)
	}
	if m.Metadata.Name == "" {
		return fmt.Errorf("metadata.name is required")
	}
	if !semverRE.MatchString(m.Metadata.Version) {
		return fmt.Errorf("metadata.version %q is not valid semver (expected MAJOR.MINOR.PATCH)", m.Metadata.Version)
	}
	if m.Metadata.TrustTier < domain.TierOfficial || m.Metadata.TrustTier > domain.TierCommunity {
		return fmt.Errorf("metadata.trustTier must be 1, 2, or 3")
	}
	if m.Metadata.CruxVersionConstraint == "" {
		return fmt.Errorf("metadata.cruxVersionConstraint is required")
	}
	return nil
}

// RunPreGenerate executes all preGenerate hooks on p in order.
func RunPreGenerate(ctx context.Context, p *domain.Plugin, hctx *domain.HookContext) error {
	for _, h := range p.PreGenerate {
		if err := h(ctx, hctx); err != nil {
			return fmt.Errorf("plugin %s preGenerate hook: %w", p.Manifest.Metadata.Name, err)
		}
	}
	return nil
}

// RunPostGenerate executes all postGenerate hooks on p in order.
func RunPostGenerate(ctx context.Context, p *domain.Plugin, hctx *domain.HookContext) error {
	for _, h := range p.PostGenerate {
		if err := h(ctx, hctx); err != nil {
			return fmt.Errorf("plugin %s postGenerate hook: %w", p.Manifest.Metadata.Name, err)
		}
	}
	return nil
}

// checkCompatibility verifies that cruxVersion satisfies the constraint string.
// Supported operators: >=, >, <=, <, = (and bare version implies =).
// Multiple constraints separated by space are ANDed.
func checkCompatibility(constraint, cruxVersion string) error {
	if constraint == "" {
		return nil
	}
	parts := strings.Fields(constraint)
	for _, part := range parts {
		ok, err := satisfies(part, cruxVersion)
		if err != nil {
			return fmt.Errorf("invalid version constraint %q: %w", part, err)
		}
		if !ok {
			return fmt.Errorf("crux version %s does not satisfy constraint %q", cruxVersion, constraint)
		}
	}
	return nil
}

// satisfies checks whether version satisfies a single constraint token like ">=1.0.0".
func satisfies(token, version string) (bool, error) {
	var op, ver string
	for _, prefix := range []string{">=", "<=", ">", "<", "="} {
		if strings.HasPrefix(token, prefix) {
			op = prefix
			ver = strings.TrimPrefix(token, prefix)
			break
		}
	}
	if op == "" {
		op = "="
		ver = token
	}

	cmp, err := compareSemver(version, ver)
	if err != nil {
		return false, err
	}
	switch op {
	case ">=":
		return cmp >= 0, nil
	case ">":
		return cmp > 0, nil
	case "<=":
		return cmp <= 0, nil
	case "<":
		return cmp < 0, nil
	default: // "="
		return cmp == 0, nil
	}
}

// compareSemver returns -1, 0, or 1 for a < b, a == b, a > b.
func compareSemver(a, b string) (int, error) {
	ap, err := parseSemver(a)
	if err != nil {
		return 0, err
	}
	bp, err := parseSemver(b)
	if err != nil {
		return 0, err
	}
	for i := range ap {
		if ap[i] < bp[i] {
			return -1, nil
		}
		if ap[i] > bp[i] {
			return 1, nil
		}
	}
	return 0, nil
}

func parseSemver(v string) ([3]int, error) {
	// Strip leading 'v' if present.
	v = strings.TrimPrefix(v, "v")
	parts := strings.SplitN(v, ".", 3)
	if len(parts) != 3 {
		return [3]int{}, fmt.Errorf("invalid semver %q", v)
	}
	var out [3]int
	for i, p := range parts {
		// Strip pre-release/build metadata suffix.
		p = strings.FieldsFunc(p, func(r rune) bool { return r == '-' || r == '+' })[0]
		n, err := strconv.Atoi(p)
		if err != nil {
			return [3]int{}, fmt.Errorf("invalid semver %q", v)
		}
		out[i] = n
	}
	return out, nil
}
