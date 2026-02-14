package compile

import (
	"fmt"
	"path"

	rcfg "github.com/kalo-build/morphe-go/pkg/registry/cfg"
	"github.com/kalo-build/plugin-morphe-ts-zod-bridge/pkg/compile/cfg"
)

// BridgeConfig holds the configuration for the TS/Zod casing bridge.
type BridgeConfig struct {
	// Registry loading configuration
	rcfg.MorpheLoadRegistryConfig

	// Output path for generated converter files
	OutputPath string

	// SourceCasing is the field casing of the source format (Zod schemas / wire format).
	SourceCasing cfg.Casing

	// TargetCasing is the field casing of the target format (TypeScript types / application format).
	TargetCasing cfg.Casing

	// ZodSchemasImportPath is the base import path for generated Zod schema files.
	ZodSchemasImportPath string

	// TsTypesImportPath is the base import path for generated TypeScript type definition files.
	TsTypesImportPath string
}

// DefaultBridgeConfig creates a default configuration.
func DefaultBridgeConfig(
	yamlRegistryPath string,
	baseOutputDirPath string,
) BridgeConfig {
	return BridgeConfig{
		MorpheLoadRegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryEnumsDirPath:      path.Join(yamlRegistryPath, "enums"),
			RegistryModelsDirPath:     path.Join(yamlRegistryPath, "models"),
			RegistryStructuresDirPath: path.Join(yamlRegistryPath, "structures"),
			RegistryEntitiesDirPath:   path.Join(yamlRegistryPath, "entities"),
		},
		OutputPath:   baseOutputDirPath,
		SourceCasing: cfg.CasingSnake,
		TargetCasing: cfg.CasingCamel,
	}
}

// Validate checks if the configuration is valid.
func (c BridgeConfig) Validate() error {
	if err := c.MorpheLoadRegistryConfig.Validate(); err != nil {
		return err
	}
	if c.OutputPath == "" {
		return fmt.Errorf("output path is required")
	}
	if !c.SourceCasing.IsValid() {
		return cfg.ErrInvalidCasing(c.SourceCasing)
	}
	if !c.TargetCasing.IsValid() {
		return cfg.ErrInvalidCasing(c.TargetCasing)
	}
	if c.ZodSchemasImportPath == "" {
		return fmt.Errorf("zodSchemasImportPath is required")
	}
	if c.TsTypesImportPath == "" {
		return fmt.Errorf("tsTypesImportPath is required")
	}
	return nil
}
