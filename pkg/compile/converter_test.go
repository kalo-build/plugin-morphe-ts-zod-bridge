package compile_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kalo-build/plugin-morphe-ts-zod-bridge/pkg/compile"
	"github.com/kalo-build/plugin-morphe-ts-zod-bridge/pkg/compile/cfg"
)

type ConverterTestSuite struct {
	suite.Suite
}

func TestConverterTestSuite(t *testing.T) {
	suite.Run(t, new(ConverterTestSuite))
}

func (s *ConverterTestSuite) TestBuildFieldMappings_SnakeToCamel() {
	// Arrange
	fields := []compile.FieldInput{
		{Name: "FirstName"},
		{Name: "LastName"},
		{Name: "DateOfBirth"},
	}

	// Act
	mappings := compile.BuildFieldMappings(fields, cfg.CasingSnake, cfg.CasingCamel)

	// Assert
	s.Require().Len(mappings, 3)
	s.Equal("first_name", mappings[0].SourceName)
	s.Equal("firstName", mappings[0].TargetName)
	s.False(mappings[0].Optional)
	s.Equal("last_name", mappings[1].SourceName)
	s.Equal("lastName", mappings[1].TargetName)
	s.Equal("date_of_birth", mappings[2].SourceName)
	s.Equal("dateOfBirth", mappings[2].TargetName)
}

func (s *ConverterTestSuite) TestBuildFieldMappings_CamelToSnake() {
	// Arrange
	fields := []compile.FieldInput{
		{Name: "Email"},
		{Name: "ZipCode"},
	}

	// Act
	mappings := compile.BuildFieldMappings(fields, cfg.CasingCamel, cfg.CasingSnake)

	// Assert
	s.Require().Len(mappings, 2)
	s.Equal("email", mappings[0].SourceName)
	s.Equal("email", mappings[0].TargetName)
	s.Equal("zipCode", mappings[1].SourceName)
	s.Equal("zip_code", mappings[1].TargetName)
}

func (s *ConverterTestSuite) TestBuildFieldMappings_SameCasing() {
	// Arrange
	fields := []compile.FieldInput{
		{Name: "FirstName"},
		{Name: "LastName"},
	}

	// Act
	mappings := compile.BuildFieldMappings(fields, cfg.CasingCamel, cfg.CasingCamel)

	// Assert
	s.Require().Len(mappings, 2)
	s.Equal("firstName", mappings[0].SourceName)
	s.Equal("firstName", mappings[0].TargetName)
	s.Equal("lastName", mappings[1].SourceName)
	s.Equal("lastName", mappings[1].TargetName)
}

func (s *ConverterTestSuite) TestBuildFieldMappings_EmptyInput() {
	// Arrange & Act
	mappings := compile.BuildFieldMappings([]compile.FieldInput{}, cfg.CasingSnake, cfg.CasingCamel)

	// Assert
	s.Empty(mappings)
}

func (s *ConverterTestSuite) TestBuildFieldMappings_OptionalFields() {
	// Arrange
	fields := []compile.FieldInput{
		{Name: "Email"},
		{Name: "Nickname", Optional: true},
	}

	// Act
	mappings := compile.BuildFieldMappings(fields, cfg.CasingSnake, cfg.CasingCamel)

	// Assert
	s.Require().Len(mappings, 2)
	s.False(mappings[0].Optional)
	s.True(mappings[1].Optional)
}

func (s *ConverterTestSuite) TestBuildFieldMappings_WithSuffix() {
	// Arrange
	fields := []compile.FieldInput{
		{Name: "PersonID", Optional: true, Suffix: "s"},
		{Name: "Person", Optional: true, Suffix: "s"},
	}

	// Act
	mappings := compile.BuildFieldMappings(fields, cfg.CasingSnake, cfg.CasingCamel)

	// Assert
	s.Require().Len(mappings, 2)
	s.Equal("person_ids", mappings[0].SourceName)
	s.Equal("personIDs", mappings[0].TargetName)
	s.True(mappings[0].Optional)
	s.Equal("persons", mappings[1].SourceName)
	s.Equal("persons", mappings[1].TargetName)
	s.True(mappings[1].Optional)
}

func (s *ConverterTestSuite) TestBuildConverterData_ModelCategory() {
	// Arrange
	fields := []compile.FieldInput{
		{Name: "FirstName"},
		{Name: "LastName"},
	}
	config := compile.BridgeConfig{
		SourceCasing:         cfg.CasingSnake,
		TargetCasing:         cfg.CasingCamel,
		ZodSchemasImportPath: "@/generated/schemas",
		TsTypesImportPath:    "@/generated/types",
	}

	// Act
	data := compile.BuildConverterData("Person", "models", fields, config)

	// Assert
	s.Equal("Person", data.TypeName)
	s.Equal("models", data.Category)
	s.Equal("convertPerson", data.FunctionName)
	s.Equal("PersonSchema", data.SchemaName)
	s.Equal("@/generated/types/models/person", data.TsTypeImportPath)
	s.Equal("@/generated/schemas/models/person", data.ZodSchemaImportPath)
	s.Require().Len(data.Fields, 2)
	s.Equal("first_name", data.Fields[0].SourceName)
	s.Equal("firstName", data.Fields[0].TargetName)
}

func (s *ConverterTestSuite) TestBuildConverterData_MultiWordTypeName() {
	// Arrange
	fields := []compile.FieldInput{
		{Name: "Street"},
		{Name: "HouseNr"},
	}
	config := compile.BridgeConfig{
		SourceCasing:         cfg.CasingSnake,
		TargetCasing:         cfg.CasingCamel,
		ZodSchemasImportPath: "@/schemas",
		TsTypesImportPath:    "@/types",
	}

	// Act
	data := compile.BuildConverterData("ContactInfo", "models", fields, config)

	// Assert
	s.Equal("ContactInfo", data.TypeName)
	s.Equal("convertContactInfo", data.FunctionName)
	s.Equal("ContactInfoSchema", data.SchemaName)
	s.Equal("@/types/models/contact-info", data.TsTypeImportPath)
	s.Equal("@/schemas/models/contact-info", data.ZodSchemaImportPath)
}

func (s *ConverterTestSuite) TestRenderConverter_SnakeToCamel() {
	// Arrange
	data := compile.ConverterData{
		TypeName:            "Person",
		Category:            "models",
		FunctionName:        "convertPerson",
		SchemaName:          "PersonSchema",
		TsTypeImportPath:    "@/generated/types/models/person",
		ZodSchemaImportPath: "@/generated/schemas/models/person",
		Fields: []compile.FieldMapping{
			{SourceName: "first_name", TargetName: "firstName"},
			{SourceName: "last_name", TargetName: "lastName"},
		},
	}

	// Act
	content := compile.RenderConverter(data)
	contentStr := string(content)

	// Assert
	s.Contains(contentStr, "// Code generated by plugin-morphe-ts-zod-bridge. DO NOT EDIT.")
	s.Contains(contentStr, "import type { Person } from '@/generated/types/models/person';")
	s.Contains(contentStr, "import { PersonSchema } from '@/generated/schemas/models/person';")
	s.Contains(contentStr, "import { z } from 'zod';")
	s.Contains(contentStr, "export function convertPerson(")
	s.Contains(contentStr, "source: z.infer<typeof PersonSchema>")
	s.Contains(contentStr, "): Person {")
	s.Contains(contentStr, "firstName: source.first_name,")
	s.Contains(contentStr, "lastName: source.last_name,")
}

func (s *ConverterTestSuite) TestRenderConverter_EmptyFields() {
	// Arrange
	data := compile.ConverterData{
		TypeName:            "Empty",
		Category:            "models",
		FunctionName:        "convertEmpty",
		SchemaName:          "EmptySchema",
		TsTypeImportPath:    "@/types/models/empty",
		ZodSchemaImportPath: "@/schemas/models/empty",
		Fields:              []compile.FieldMapping{},
	}

	// Act
	content := compile.RenderConverter(data)
	contentStr := string(content)

	// Assert
	s.Contains(contentStr, "export function convertEmpty(")
	s.Contains(contentStr, "return {\n  };\n}")
}

func (s *ConverterTestSuite) TestRenderConverter_SingleField() {
	// Arrange
	data := compile.ConverterData{
		TypeName:            "Simple",
		Category:            "structures",
		FunctionName:        "convertSimple",
		SchemaName:          "SimpleSchema",
		TsTypeImportPath:    "@/types/structures/simple",
		ZodSchemaImportPath: "@/schemas/structures/simple",
		Fields: []compile.FieldMapping{
			{SourceName: "name", TargetName: "name"},
		},
	}

	// Act
	content := compile.RenderConverter(data)
	contentStr := string(content)

	// Assert
	s.Contains(contentStr, "name: source.name,")
	s.Contains(contentStr, "): Simple {")
}

func (s *ConverterTestSuite) TestRenderConverter_OptionalFields() {
	// Arrange
	data := compile.ConverterData{
		TypeName:            "User",
		Category:            "models",
		FunctionName:        "convertUser",
		SchemaName:          "UserSchema",
		TsTypeImportPath:    "@/types/models/user",
		ZodSchemaImportPath: "@/schemas/models/user",
		Fields: []compile.FieldMapping{
			{SourceName: "email", TargetName: "email", Optional: false},
			{SourceName: "nickname", TargetName: "nickname", Optional: true},
			{SourceName: "organization_id", TargetName: "organizationID", Optional: true},
		},
	}

	// Act
	content := compile.RenderConverter(data)
	contentStr := string(content)

	// Assert: required field uses direct assignment
	s.Contains(contentStr, "    email: source.email,\n")
	// Assert: optional fields use conditional spread
	s.Contains(contentStr, "    ...(source.nickname !== undefined && {\n")
	s.Contains(contentStr, "      nickname: source.nickname,\n")
	s.Contains(contentStr, "    }),\n")
	s.Contains(contentStr, "    ...(source.organization_id !== undefined && {\n")
	s.Contains(contentStr, "      organizationID: source.organization_id,\n")
}

func (s *ConverterTestSuite) TestRenderConverter_RelationshipFields() {
	// Arrange
	data := compile.ConverterData{
		TypeName:            "Person",
		Category:            "models",
		FunctionName:        "convertPerson",
		SchemaName:          "PersonSchema",
		TsTypeImportPath:    "@/types/models/person",
		ZodSchemaImportPath: "@/schemas/models/person",
		Fields: []compile.FieldMapping{
			{SourceName: "id", TargetName: "id"},
			{SourceName: "company_id", TargetName: "companyID", Optional: true},
			{SourceName: "company", TargetName: "company", Optional: true, ConverterFn: "convertCompany", ConverterImport: "./company"},
			{SourceName: "person_ids", TargetName: "personIDs", Optional: true},
			{SourceName: "persons", TargetName: "persons", Optional: true, ConverterFn: "convertPerson", ConverterImport: "./person", IsArray: true},
		},
	}

	// Act
	content := compile.RenderConverter(data)
	contentStr := string(content)

	// Assert: converter import for Company (not self-reference)
	s.Contains(contentStr, "import { convertCompany } from './company';")
	// Assert: no import for convertPerson (self-reference)
	s.NotContains(contentStr, "import { convertPerson }")

	// Assert: single relationship uses converter call
	s.Contains(contentStr, "company: convertCompany(source.company),")
	// Assert: array relationship uses .map()
	s.Contains(contentStr, "persons: source.persons.map(convertPerson),")

	// Assert: FK ID fields remain as-is (no converter call)
	s.Contains(contentStr, "companyID: source.company_id,")
	s.Contains(contentStr, "personIDs: source.person_ids,")
}

func (s *ConverterTestSuite) TestBuildFieldMappings_WithRelationTarget() {
	// Arrange
	fields := []compile.FieldInput{
		{Name: "CompanyID", Optional: true},
		{Name: "Company", Optional: true, RelationTarget: "Company"},
		{Name: "PersonID", Optional: true, Suffix: "s"},
		{Name: "Person", Optional: true, Suffix: "s", RelationTarget: "Person", IsArray: true},
	}

	// Act
	mappings := compile.BuildFieldMappings(fields, cfg.CasingSnake, cfg.CasingCamel)

	// Assert
	s.Require().Len(mappings, 4)

	// FK ID field: no converter
	s.Equal("company_id", mappings[0].SourceName)
	s.Empty(mappings[0].ConverterFn)

	// Single relationship object: converter function + import
	s.Equal("company", mappings[1].SourceName)
	s.Equal("convertCompany", mappings[1].ConverterFn)
	s.Equal("./company", mappings[1].ConverterImport)
	s.False(mappings[1].IsArray)

	// Array FK IDs: no converter
	s.Equal("person_ids", mappings[2].SourceName)
	s.Empty(mappings[2].ConverterFn)

	// Array relationship object: converter function + import + IsArray
	s.Equal("persons", mappings[3].SourceName)
	s.Equal("convertPerson", mappings[3].ConverterFn)
	s.Equal("./person", mappings[3].ConverterImport)
	s.True(mappings[3].IsArray)
}

func (s *ConverterTestSuite) TestRenderConverter_MixedRequiredAndOptional() {
	// Arrange
	data := compile.ConverterData{
		TypeName:            "Profile",
		Category:            "structures",
		FunctionName:        "convertProfile",
		SchemaName:          "ProfileSchema",
		TsTypeImportPath:    "@/types/structures/profile",
		ZodSchemaImportPath: "@/schemas/structures/profile",
		Fields: []compile.FieldMapping{
			{SourceName: "display_name", TargetName: "displayName", Optional: false},
			{SourceName: "avatar_url", TargetName: "avatarURL", Optional: true},
			{SourceName: "email", TargetName: "email", Optional: false},
		},
	}

	// Act
	content := compile.RenderConverter(data)
	contentStr := string(content)

	// Assert: required fields (target is camel, source is snake)
	s.Contains(contentStr, "    displayName: source.display_name,\n")
	s.Contains(contentStr, "    email: source.email,\n")
	// Assert: optional field
	s.Contains(contentStr, "    ...(source.avatar_url !== undefined && {\n")
}
