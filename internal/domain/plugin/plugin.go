// Package plugin defines the plugin domain types, manifest schema, and loader interface.
package plugin

import "context"

// TrustTier classifies a plugin's provenance and review level.
type TrustTier int

const (
	TierOfficial  TrustTier = 1 // Platform-team owned, bundled with crux
	TierVerified  TrustTier = 2 // Inner-source, platform-reviewed
	TierCommunity TrustTier = 3 // Self-published, unvetted
)

// HookContext carries the data available to lifecycle hooks.
type HookContext struct {
	ServiceName string
	OutputDir   string
	Answers     map[string]any
}

// Hook is a function invoked at a lifecycle point during generation.
type Hook func(ctx context.Context, hctx *HookContext) error

// QuestionSpec is a single prompt contributed by a plugin.
type QuestionSpec struct {
	ID      string   `yaml:"id"`
	Type    string   `yaml:"type"`
	Prompt  string   `yaml:"prompt"`
	Options []string `yaml:"options,omitempty"`
	Default string   `yaml:"default,omitempty"`
}

// HooksSpec lists hook command strings declared in the manifest.
// At MVP these are resolved to registered Go Hook functions by name.
type HooksSpec struct {
	PreGenerate  []string `yaml:"preGenerate,omitempty"`
	PostGenerate []string `yaml:"postGenerate,omitempty"`
}

// Metadata holds the plugin identity fields from plugin.yaml.
type Metadata struct {
	Name                  string    `yaml:"name"`
	Version               string    `yaml:"version"`
	Description           string    `yaml:"description"`
	Author                string    `yaml:"author"`
	TrustTier             TrustTier `yaml:"trustTier"`
	Tags                  []string  `yaml:"tags,omitempty"`
	CruxVersionConstraint string    `yaml:"cruxVersionConstraint"`
}

// Spec holds the capability declarations from plugin.yaml.
type Spec struct {
	Questions    []QuestionSpec `yaml:"questions,omitempty"`
	Templates    []string       `yaml:"templates,omitempty"`
	Hooks        HooksSpec      `yaml:"hooks,omitempty"`
	Dependencies []string       `yaml:"dependencies,omitempty"`
}

// Manifest is the parsed, validated representation of a plugin.yaml file.
type Manifest struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       Spec     `yaml:"spec"`
}

// Plugin is a loaded, ready-to-use plugin with its manifest and registered hooks.
type Plugin struct {
	Manifest     *Manifest
	PreGenerate  []Hook
	PostGenerate []Hook
}

// Loader discovers and loads plugins, returning them ready for use.
type Loader interface {
	// Load discovers all plugins from the embedded dir and ~/.crux/plugins/,
	// validates each manifest against cruxVersion, and returns loaded plugins.
	Load(cruxVersion string) ([]*Plugin, error)
}
