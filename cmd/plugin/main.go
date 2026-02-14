package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kalo-build/plugin-morphe-ts-zod-bridge/pkg/compile"
	"github.com/kalo-build/plugin-morphe-ts-zod-bridge/pkg/compile/cfg"
)

// PluginConfig represents the configuration passed to the plugin by Kalo CLI.
type PluginConfig struct {
	InputPath  string                 `json:"inputPath"`
	OutputPath string                 `json:"outputPath"`
	Config     map[string]interface{} `json:"config,omitempty"`
	Verbose    bool                   `json:"verbose,omitempty"`
}

// Exit codes
const (
	ExitSuccess         = 0
	ExitCompileFailed   = 1
	ExitMissingConfig   = 3
	ExitInvalidConfig   = 4
	ExitInputPathError  = 12
	ExitOutputPathError = 13
)

// logInfo prints info messages only when verbose mode is enabled.
func logInfo(verbose bool, format string, args ...interface{}) {
	if verbose {
		fmt.Fprintf(os.Stdout, format+"\n", args...)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: plugin-morphe-ts-zod-bridge <config>")
		fmt.Fprintln(os.Stderr, "  config: JSON string with inputPath, outputPath, and config parameters")
		os.Exit(ExitMissingConfig)
	}

	// Parse configuration
	rawConfig := os.Args[1]
	var pluginConfig PluginConfig
	if err := json.Unmarshal([]byte(rawConfig), &pluginConfig); err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing config JSON:", err)
		os.Exit(ExitInvalidConfig)
	}

	// Validate required paths
	if pluginConfig.InputPath == "" {
		fmt.Fprintln(os.Stderr, "Error: inputPath is required")
		os.Exit(ExitInputPathError)
	}
	if pluginConfig.OutputPath == "" {
		fmt.Fprintln(os.Stderr, "Error: outputPath is required")
		os.Exit(ExitOutputPathError)
	}

	// Convert to absolute paths
	if absPath, err := filepath.Abs(pluginConfig.InputPath); err == nil {
		pluginConfig.InputPath = absPath
	}
	if absPath, err := filepath.Abs(pluginConfig.OutputPath); err == nil {
		pluginConfig.OutputPath = absPath
	}

	// Extract plugin-specific config
	sourceCasing := getStringConfig(pluginConfig.Config, "sourceCasing")
	targetCasing := getStringConfig(pluginConfig.Config, "targetCasing")
	zodSchemasImportPath := getStringConfig(pluginConfig.Config, "zodSchemasImportPath")
	tsTypesImportPath := getStringConfig(pluginConfig.Config, "tsTypesImportPath")

	// Apply defaults
	if sourceCasing == "" {
		sourceCasing = "snake"
	}
	if targetCasing == "" {
		targetCasing = "camel"
	}

	if zodSchemasImportPath == "" {
		fmt.Fprintln(os.Stderr, "Error: zodSchemasImportPath is required in config")
		os.Exit(ExitInvalidConfig)
	}
	if tsTypesImportPath == "" {
		fmt.Fprintln(os.Stderr, "Error: tsTypesImportPath is required in config")
		os.Exit(ExitInvalidConfig)
	}

	logInfo(pluginConfig.Verbose, "Registry path: %s", pluginConfig.InputPath)
	logInfo(pluginConfig.Verbose, "Output path: %s", pluginConfig.OutputPath)
	logInfo(pluginConfig.Verbose, "Source casing: %s", sourceCasing)
	logInfo(pluginConfig.Verbose, "Target casing: %s", targetCasing)

	// Build compile configuration
	bridgeConfig := compile.DefaultBridgeConfig(pluginConfig.InputPath, pluginConfig.OutputPath)
	bridgeConfig.SourceCasing = cfg.Casing(sourceCasing)
	bridgeConfig.TargetCasing = cfg.Casing(targetCasing)
	bridgeConfig.ZodSchemasImportPath = zodSchemasImportPath
	bridgeConfig.TsTypesImportPath = tsTypesImportPath

	if err := bridgeConfig.Validate(); err != nil {
		fmt.Fprintln(os.Stderr, "Invalid configuration:", err)
		os.Exit(ExitInvalidConfig)
	}

	logInfo(pluginConfig.Verbose, "Starting TS/Zod bridge generation...")
	if err := compile.MorpheToBridge(bridgeConfig); err != nil {
		fmt.Fprintln(os.Stderr, "Bridge generation failed:", err)
		os.Exit(ExitCompileFailed)
	}

	logInfo(pluginConfig.Verbose, "Bridge generation completed successfully")
	os.Exit(ExitSuccess)
}

func getStringConfig(config map[string]interface{}, key string) string {
	if config == nil {
		return ""
	}
	if v, ok := config[key].(string); ok {
		return v
	}
	return ""
}
