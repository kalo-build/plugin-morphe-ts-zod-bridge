package compile

import (
	"fmt"

	"github.com/kalo-build/morphe-go/pkg/registry"
)

// MorpheToBridge generates TypeScript converter functions that bridge between
// Zod schema casing and TypeScript type casing for all types in the Morphe registry.
func MorpheToBridge(config BridgeConfig) error {
	// Load the Morphe registry
	r, rErr := registry.LoadMorpheRegistry(registry.LoadMorpheRegistryHooks{}, config.MorpheLoadRegistryConfig)
	if rErr != nil {
		return fmt.Errorf("failed to load morphe registry: %w", rErr)
	}

	writer := NewBridgeWriter(config.OutputPath)

	// Compile models
	if r.HasModels() {
		if err := CompileAllModelConverters(config, r, writer); err != nil {
			return fmt.Errorf("failed to compile model converters: %w", err)
		}
	}

	// Compile structures
	if r.HasStructures() {
		if err := CompileAllStructureConverters(config, r, writer); err != nil {
			return fmt.Errorf("failed to compile structure converters: %w", err)
		}
	}

	// Compile entities
	if r.HasEntities() {
		if err := CompileAllEntityConverters(config, r, writer); err != nil {
			return fmt.Errorf("failed to compile entity converters: %w", err)
		}
	}

	return nil
}
