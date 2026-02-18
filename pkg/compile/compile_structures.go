package compile

import (
	"fmt"

	"github.com/kalo-build/go-util/core"
	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
)

// CompileAllStructureConverters generates converter functions for all structures in the registry.
func CompileAllStructureConverters(config BridgeConfig, r *registry.Registry, writer *BridgeWriter) error {
	for structureName, structure := range r.GetAllStructures() {
		fields := collectStructureFieldInputs(structure)

		data := BuildConverterData(structureName, "structures", fields, config)
		content := RenderConverter(data)

		if err := writer.WriteStructureConverter(structureName, content); err != nil {
			return fmt.Errorf("failed to write structure converter %q: %w", structureName, err)
		}
	}
	return nil
}

// collectStructureFieldInputs builds FieldInput entries for direct structure fields,
// checking each field's Attributes for the "optional" marker.
func collectStructureFieldInputs(structure yaml.Structure) []FieldInput {
	fieldNames := core.MapKeysSorted(structure.Fields)
	inputs := make([]FieldInput, 0, len(fieldNames))
	for _, name := range fieldNames {
		field := structure.Fields[name]
		inputs = append(inputs, FieldInput{
			Name:     name,
			Optional: hasAttribute(field.Attributes, "optional"),
		})
	}
	return inputs
}
