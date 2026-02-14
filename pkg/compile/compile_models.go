package compile

import (
	"fmt"

	"github.com/kalo-build/go-util/core"
	"github.com/kalo-build/morphe-go/pkg/registry"
)

// CompileAllModelConverters generates converter functions for all models in the registry.
func CompileAllModelConverters(config BridgeConfig, r *registry.Registry, writer *BridgeWriter) error {
	for modelName, model := range r.GetAllModels() {
		// Collect direct field names (sorted for deterministic output)
		fieldNames := core.MapKeysSorted(model.Fields)

		data := BuildConverterData(modelName, "models", fieldNames, config)
		content := RenderConverter(data)

		if err := writer.WriteModelConverter(modelName, content); err != nil {
			return fmt.Errorf("failed to write model converter %q: %w", modelName, err)
		}
	}
	return nil
}
