// Package config holds application-level configuration shared across commands.
package config

// GlobalConfig holds flags available on every command.
type GlobalConfig struct {
	Verbose    bool
	OutputMode string // "text" or "json"
	ConfigFile string
}
