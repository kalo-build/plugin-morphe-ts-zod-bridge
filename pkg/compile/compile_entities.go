package compile

import (
	"fmt"
	"strings"

	"github.com/kalo-build/go-util/core"
	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/morphe-go/pkg/yamlops"
)

// CompileAllEntityConverters generates converter functions for all entities in the registry.
func CompileAllEntityConverters(config BridgeConfig, r *registry.Registry, writer *BridgeWriter) error {
	for entityName, entity := range r.GetAllEntities() {
		fields := collectEntityFieldInputs(entity)

		relationFields, relErr := collectEntityRelationFieldInputs(r, entity.Related)
		if relErr != nil {
			return fmt.Errorf("entity %q: %w", entityName, relErr)
		}
		fields = append(fields, relationFields...)

		data := BuildConverterData(entityName, "entities", fields, config)
		content := RenderConverter(data)

		if err := writer.WriteEntityConverter(entityName, content); err != nil {
			return fmt.Errorf("failed to write entity converter %q: %w", entityName, err)
		}
	}
	return nil
}

// collectEntityFieldInputs builds FieldInput entries for direct entity fields,
// checking each field's Attributes for the "optional" marker.
func collectEntityFieldInputs(entity yaml.Entity) []FieldInput {
	fieldNames := core.MapKeysSorted(entity.Fields)
	inputs := make([]FieldInput, 0, len(fieldNames))
	for _, name := range fieldNames {
		field := entity.Fields[name]
		inputs = append(inputs, FieldInput{
			Name:     name,
			Optional: hasAttribute(field.Attributes, "optional"),
		})
	}
	return inputs
}

// collectEntityRelationFieldInputs builds FieldInput entries for entity
// relationship-derived fields. Entity relations reference other entities, but
// the FK ID is derived from the underlying model's primary identifier. FK/ID
// fields follow the relation's optionality attribute; relation object fields
// are always optional per ADR-003.
func collectEntityRelationFieldInputs(r *registry.Registry, relations map[string]yaml.EntityRelation) ([]FieldInput, error) {
	if len(relations) == 0 {
		return nil, nil
	}

	var inputs []FieldInput
	relNames := core.MapKeysSorted(relations)
	for _, relName := range relNames {
		rel := relations[relName]

		if yamlops.IsRelationPolyFor(rel.Type) {
			inputs = append(inputs, collectForPolyFieldInputs(relName, rel.Attributes)...)
			continue
		}

		targetEntityName := relName
		if rel.Aliased != "" {
			targetEntityName = rel.Aliased
		} else if yamlops.IsRelationPolyHas(rel.Type) && len(rel.For) > 0 {
			targetEntityName = rel.For[0]
		}

		targetEntity, entErr := r.GetEntity(targetEntityName)
		if entErr != nil {
			return nil, fmt.Errorf("relationship %q: %w", relName, entErr)
		}

		primaryIDFieldName, idErr := getEntityPrimaryIdentifierFieldName(targetEntity)
		if idErr != nil {
			return nil, fmt.Errorf("relationship %q: %w", relName, idErr)
		}

		isMany := yamlops.IsRelationMany(rel.Type)
		suffix := ""
		if isMany {
			suffix = "s"
		}

		fkOptional := !isMany && hasAttribute(rel.Attributes, "optional")
		inputs = append(inputs, FieldInput{
			Name:     relName + primaryIDFieldName,
			Optional: fkOptional,
			Suffix:   suffix,
		})
		inputs = append(inputs, FieldInput{
			Name:           relName,
			Optional:       true,
			Suffix:         suffix,
			RelationTarget: targetEntityName,
			IsArray:        isMany,
		})
	}
	return inputs, nil
}

// getEntityPrimaryIdentifierFieldName extracts the primary identifier field
// name from an entity definition. Entities follow the same identifier convention
// as models: identifiers.primary is a single field name.
func getEntityPrimaryIdentifierFieldName(entity yaml.Entity) (string, error) {
	primaryID, hasPrimaryID := entity.Identifiers["primary"]
	if !hasPrimaryID {
		return "", fmt.Errorf("entity %q has no primary identifier", entity.Name)
	}
	if len(primaryID.Fields) != 1 {
		return "", fmt.Errorf("entity %q primary identifier must have exactly one field", entity.Name)
	}
	field := primaryID.Fields[0]
	if strings.HasPrefix(field, "rel:") {
		field = strings.TrimPrefix(field, "rel:") + "ID"
	}
	return field, nil
}
