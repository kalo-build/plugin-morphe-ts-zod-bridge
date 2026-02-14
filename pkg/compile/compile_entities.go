package compile

import (
	"fmt"

	"github.com/kalo-build/go-util/core"
	"github.com/kalo-build/morphe-go/pkg/registry"
)

// CompileAllEntityConverters generates converter functions for all entities in the registry.
func CompileAllEntityConverters(config BridgeConfig, r *registry.Registry, writer *BridgeWriter) error {
	for entityName, entity := range r.GetAllEntities() {
		// Collect direct field names (sorted for deterministic output)
		fieldNames := core.MapKeysSorted(entity.Fields)

		data := BuildConverterData(entityName, "entities", fieldNames, config)
		content := RenderConverter(data)

		if err := writer.WriteEntityConverter(entityName, content); err != nil {
			return fmt.Errorf("failed to write entity converter %q: %w", entityName, err)
		}
	}
	return nil
}
