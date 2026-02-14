package compile

import (
	"fmt"

	"github.com/kalo-build/go-util/core"
	"github.com/kalo-build/morphe-go/pkg/registry"
)

// CompileAllStructureConverters generates converter functions for all structures in the registry.
func CompileAllStructureConverters(config BridgeConfig, r *registry.Registry, writer *BridgeWriter) error {
	for structureName, structure := range r.GetAllStructures() {
		// Collect direct field names (sorted for deterministic output)
		fieldNames := core.MapKeysSorted(structure.Fields)

		data := BuildConverterData(structureName, "structures", fieldNames, config)
		content := RenderConverter(data)

		if err := writer.WriteStructureConverter(structureName, content); err != nil {
			return fmt.Errorf("failed to write structure converter %q: %w", structureName, err)
		}
	}
	return nil
}
