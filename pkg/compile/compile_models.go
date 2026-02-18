package compile

import (
	"fmt"

	"github.com/kalo-build/go-util/core"
	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/morphe-go/pkg/yamlops"
)

// CompileAllModelConverters generates converter functions for all models in the registry.
func CompileAllModelConverters(config BridgeConfig, r *registry.Registry, writer *BridgeWriter) error {
	for modelName, model := range r.GetAllModels() {
		fields := collectModelFieldInputs(model)

		relationFields, relErr := collectModelRelationFieldInputs(r, model.Related)
		if relErr != nil {
			return fmt.Errorf("model %q: %w", modelName, relErr)
		}
		fields = append(fields, relationFields...)

		data := BuildConverterData(modelName, "models", fields, config)
		content := RenderConverter(data)

		if err := writer.WriteModelConverter(modelName, content); err != nil {
			return fmt.Errorf("failed to write model converter %q: %w", modelName, err)
		}
	}
	return nil
}

// collectModelFieldInputs builds FieldInput entries for direct model fields,
// checking each field's Attributes for the "optional" marker.
func collectModelFieldInputs(model yaml.Model) []FieldInput {
	fieldNames := core.MapKeysSorted(model.Fields)
	inputs := make([]FieldInput, 0, len(fieldNames))
	for _, name := range fieldNames {
		field := model.Fields[name]
		inputs = append(inputs, FieldInput{
			Name:     name,
			Optional: hasAttribute(field.Attributes, "optional"),
		})
	}
	return inputs
}

// collectModelRelationFieldInputs builds FieldInput entries for relationship-
// derived fields (FK ID + optional object reference). All relationship fields
// are always optional.
func collectModelRelationFieldInputs(r *registry.Registry, relations map[string]yaml.ModelRelation) ([]FieldInput, error) {
	var inputs []FieldInput

	relNames := core.MapKeysSorted(relations)
	for _, relName := range relNames {
		rel := relations[relName]

		if yamlops.IsRelationPolyFor(rel.Type) {
			inputs = append(inputs, collectForPolyFieldInputs(relName)...)
			continue
		}

		targetModelName := relName
		if rel.Aliased != "" {
			targetModelName = rel.Aliased
		} else if yamlops.IsRelationPolyHas(rel.Type) && len(rel.For) > 0 {
			targetModelName = rel.For[0]
		}

		targetModel, err := r.GetModel(targetModelName)
		if err != nil {
			return nil, fmt.Errorf("relationship %q: %w", relName, err)
		}

		primaryIDFieldName, idErr := yamlops.GetModelPrimaryIdentifierFieldName(targetModel)
		if idErr != nil {
			return nil, fmt.Errorf("relationship %q: %w", relName, idErr)
		}

		isMany := yamlops.IsRelationMany(rel.Type)
		suffix := ""
		if isMany {
			suffix = "s"
		}

		inputs = append(inputs, FieldInput{
			Name:     relName + primaryIDFieldName,
			Optional: true,
			Suffix:   suffix,
		})
		inputs = append(inputs, FieldInput{
			Name:           relName,
			Optional:       true,
			Suffix:         suffix,
			RelationTarget: targetModelName,
			IsArray:        isMany,
		})
	}
	return inputs, nil
}

// collectForPolyFieldInputs builds FieldInput entries for ForOnePoly/ForManyPoly
// relationships. These produce discriminator fields ({rel}Type + {rel}ID) with
// hardcoded suffixes — no target model lookup is needed because the ID and Type
// suffixes are fixed by convention.
func collectForPolyFieldInputs(relName string) []FieldInput {
	return []FieldInput{
		{Name: relName + "ID", Optional: true},
		{Name: relName + "Type", Optional: true},
	}
}
