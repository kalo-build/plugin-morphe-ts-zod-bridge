package compile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"

	rcfg "github.com/kalo-build/morphe-go/pkg/registry/cfg"
	"github.com/kalo-build/plugin-morphe-ts-zod-bridge/internal/testutils"
	"github.com/kalo-build/plugin-morphe-ts-zod-bridge/pkg/compile"
	"github.com/kalo-build/plugin-morphe-ts-zod-bridge/pkg/compile/cfg"
)

type ModelConverterTestSuite struct {
	suite.Suite

	TestDirPath       string
	ModelsDirPath     string
	EnumsDirPath      string
	StructuresDirPath string
	EntitiesDirPath   string
}

func TestModelConverterTestSuite(t *testing.T) {
	suite.Run(t, new(ModelConverterTestSuite))
}

func (s *ModelConverterTestSuite) SetupTest() {
	s.TestDirPath = testutils.GetTestDirPath()
	s.ModelsDirPath = filepath.Join(s.TestDirPath, "registry", "minimal", "models")
	s.EnumsDirPath = filepath.Join(s.TestDirPath, "registry", "minimal", "enums")
	s.StructuresDirPath = filepath.Join(s.TestDirPath, "registry", "minimal", "structures")
	s.EntitiesDirPath = filepath.Join(s.TestDirPath, "registry", "minimal", "entities")
}

func (s *ModelConverterTestSuite) TestCompileAllModelConverters_SnakeToCamel() {
	// Arrange
	outputDir := filepath.Join(s.TestDirPath, "working-models-snake-camel")
	s.Nil(os.Mkdir(outputDir, 0755))
	defer os.RemoveAll(outputDir)

	config := compile.BridgeConfig{
		MorpheLoadRegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryModelsDirPath:     s.ModelsDirPath,
			RegistryEnumsDirPath:      s.EnumsDirPath,
			RegistryStructuresDirPath: s.StructuresDirPath,
			RegistryEntitiesDirPath:   s.EntitiesDirPath,
		},
		OutputPath:           outputDir,
		SourceCasing:         cfg.CasingSnake,
		TargetCasing:         cfg.CasingCamel,
		ZodSchemasImportPath: "@/generated/schemas",
		TsTypesImportPath:    "@/generated/types",
	}

	// Act
	err := compile.MorpheToBridge(config)

	// Assert
	s.NoError(err)

	// Verify Person model converter
	personPath := filepath.Join(outputDir, "models", "person.ts")
	s.FileExists(personPath)
	personContent, readErr := os.ReadFile(personPath)
	s.NoError(readErr)
	personStr := string(personContent)

	s.Contains(personStr, "import type { Person } from '@/generated/types/models/person';")
	s.Contains(personStr, "import { PersonSchema } from '@/generated/schemas/models/person';")
	s.Contains(personStr, "export function convertPerson(")
	s.Contains(personStr, "firstName: source.first_name,")
	s.Contains(personStr, "lastName: source.last_name,")
	s.Contains(personStr, "nationality: source.nationality,")
	s.Contains(personStr, "id: source.id,")

	// Verify Company model converter
	companyPath := filepath.Join(outputDir, "models", "company.ts")
	s.FileExists(companyPath)
	companyContent, readErr := os.ReadFile(companyPath)
	s.NoError(readErr)
	companyStr := string(companyContent)

	s.Contains(companyStr, "import type { Company } from '@/generated/types/models/company';")
	s.Contains(companyStr, "export function convertCompany(")
	s.Contains(companyStr, "name: source.name,")
	s.Contains(companyStr, "taxID: source.tax_id,")

	// Verify ContactInfo model converter
	contactPath := filepath.Join(outputDir, "models", "contact-info.ts")
	s.FileExists(contactPath)
	contactContent, readErr := os.ReadFile(contactPath)
	s.NoError(readErr)
	contactStr := string(contactContent)

	s.Contains(contactStr, "export function convertContactInfo(")
	s.Contains(contactStr, "email: source.email,")
}

func (s *ModelConverterTestSuite) TestCompileAllModelConverters_CamelToSnake() {
	// Arrange
	outputDir := filepath.Join(s.TestDirPath, "working-models-camel-snake")
	s.Nil(os.Mkdir(outputDir, 0755))
	defer os.RemoveAll(outputDir)

	config := compile.BridgeConfig{
		MorpheLoadRegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryModelsDirPath:     s.ModelsDirPath,
			RegistryEnumsDirPath:      s.EnumsDirPath,
			RegistryStructuresDirPath: s.StructuresDirPath,
			RegistryEntitiesDirPath:   s.EntitiesDirPath,
		},
		OutputPath:           outputDir,
		SourceCasing:         cfg.CasingCamel,
		TargetCasing:         cfg.CasingSnake,
		ZodSchemasImportPath: "@/schemas",
		TsTypesImportPath:    "@/types",
	}

	// Act
	err := compile.MorpheToBridge(config)

	// Assert
	s.NoError(err)

	personPath := filepath.Join(outputDir, "models", "person.ts")
	s.FileExists(personPath)
	personContent, readErr := os.ReadFile(personPath)
	s.NoError(readErr)
	personStr := string(personContent)

	// Source is camelCase, target is snake_case
	s.Contains(personStr, "first_name: source.firstName,")
	s.Contains(personStr, "last_name: source.lastName,")
}

func (s *ModelConverterTestSuite) TestCompileAllModelConverters_RelationshipFields() {
	// Arrange
	outputDir := filepath.Join(s.TestDirPath, "working-models-relationships")
	s.Nil(os.Mkdir(outputDir, 0755))
	defer os.RemoveAll(outputDir)

	config := compile.BridgeConfig{
		MorpheLoadRegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryModelsDirPath:     s.ModelsDirPath,
			RegistryEnumsDirPath:      s.EnumsDirPath,
			RegistryStructuresDirPath: s.StructuresDirPath,
			RegistryEntitiesDirPath:   s.EntitiesDirPath,
		},
		OutputPath:           outputDir,
		SourceCasing:         cfg.CasingSnake,
		TargetCasing:         cfg.CasingCamel,
		ZodSchemasImportPath: "@/generated/schemas",
		TsTypesImportPath:    "@/generated/types",
	}

	// Act
	err := compile.MorpheToBridge(config)
	s.NoError(err)

	// Verify Person model converter includes relationship fields
	personPath := filepath.Join(outputDir, "models", "person.ts")
	personContent, readErr := os.ReadFile(personPath)
	s.NoError(readErr)
	personStr := string(personContent)

	// Person has ForOne Company and HasOne ContactInfo - both produce optional FK + object fields
	s.Contains(personStr, "import { convertCompany } from './company';")
	s.Contains(personStr, "import { convertContactInfo } from './contact-info';")
	s.Contains(personStr, "...(source.company_id !== undefined && {")
	s.Contains(personStr, "companyID: source.company_id,")
	s.Contains(personStr, "...(source.company !== undefined && {")
	s.Contains(personStr, "company: convertCompany(source.company),")
	s.Contains(personStr, "...(source.contact_info_id !== undefined && {")
	s.Contains(personStr, "contactInfoID: source.contact_info_id,")
	s.Contains(personStr, "...(source.contact_info !== undefined && {")
	s.Contains(personStr, "contactInfo: convertContactInfo(source.contact_info),")

	// Verify Company model converter includes HasMany Person relationship (pluralized)
	companyPath := filepath.Join(outputDir, "models", "company.ts")
	companyContent, readErr := os.ReadFile(companyPath)
	s.NoError(readErr)
	companyStr := string(companyContent)

	s.Contains(companyStr, "import { convertPerson } from './person';")
	s.Contains(companyStr, "...(source.person_ids !== undefined && {")
	s.Contains(companyStr, "personIDs: source.person_ids,")
	s.Contains(companyStr, "...(source.persons !== undefined && {")
	s.Contains(companyStr, "persons: source.persons.map(convertPerson),")
}

func (s *ModelConverterTestSuite) TestCompileAllModelConverters_PolymorphicRelationships() {
	// Arrange — uses the polymorphic testdata registry
	polyModelsDirPath := filepath.Join(s.TestDirPath, "registry", "polymorphic", "models")
	polyEnumsDirPath := filepath.Join(s.TestDirPath, "registry", "polymorphic", "enums")
	polyStructuresDirPath := filepath.Join(s.TestDirPath, "registry", "polymorphic", "structures")
	polyEntitiesDirPath := filepath.Join(s.TestDirPath, "registry", "polymorphic", "entities")

	outputDir := filepath.Join(s.TestDirPath, "working-models-poly")
	s.Nil(os.Mkdir(outputDir, 0755))
	defer os.RemoveAll(outputDir)

	config := compile.BridgeConfig{
		MorpheLoadRegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryModelsDirPath:     polyModelsDirPath,
			RegistryEnumsDirPath:      polyEnumsDirPath,
			RegistryStructuresDirPath: polyStructuresDirPath,
			RegistryEntitiesDirPath:   polyEntitiesDirPath,
		},
		OutputPath:           outputDir,
		SourceCasing:         cfg.CasingSnake,
		TargetCasing:         cfg.CasingCamel,
		ZodSchemasImportPath: "@/generated/schemas",
		TsTypesImportPath:    "@/generated/types",
	}

	// Act
	err := compile.MorpheToBridge(config)
	s.NoError(err)

	// Assert — Comment model has ForOnePoly Commentable [Person, Company]
	commentPath := filepath.Join(outputDir, "models", "comment.ts")
	s.FileExists(commentPath)
	commentContent, readErr := os.ReadFile(commentPath)
	s.NoError(readErr)
	commentStr := string(commentContent)

	// ForOnePoly produces discriminator fields: CommentableID and CommentableType
	s.Contains(commentStr, "...(source.commentable_id !== undefined && {")
	s.Contains(commentStr, "commentableID: source.commentable_id,")
	s.Contains(commentStr, "...(source.commentable_type !== undefined && {")
	s.Contains(commentStr, "commentableType: source.commentable_type,")

	// Assert — Person model has HasOnePoly Note (Aliased: Comment)
	personPath := filepath.Join(outputDir, "models", "person.ts")
	s.FileExists(personPath)
	personContent, readErr := os.ReadFile(personPath)
	s.NoError(readErr)
	personStr := string(personContent)

	// HasOnePoly produces FK + object reference (same as HasOne, resolved via aliased model)
	s.Contains(personStr, "import { convertComment } from './comment';")
	s.Contains(personStr, "...(source.note_id !== undefined && {")
	s.Contains(personStr, "noteID: source.note_id,")
	s.Contains(personStr, "...(source.note !== undefined && {")
	s.Contains(personStr, "note: convertComment(source.note),")
}

func (s *ModelConverterTestSuite) TestCompileAllModelConverters_OptionalFields() {
	// Arrange -- uses the optional testdata registry
	optModelsDirPath := filepath.Join(s.TestDirPath, "registry", "optional", "models")
	optEnumsDirPath := filepath.Join(s.TestDirPath, "registry", "optional", "enums")
	optStructuresDirPath := filepath.Join(s.TestDirPath, "registry", "optional", "structures")
	optEntitiesDirPath := filepath.Join(s.TestDirPath, "registry", "optional", "entities")

	outputDir := filepath.Join(s.TestDirPath, "working-models-optional")
	s.Nil(os.Mkdir(outputDir, 0755))
	defer os.RemoveAll(outputDir)

	config := compile.BridgeConfig{
		MorpheLoadRegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryModelsDirPath:     optModelsDirPath,
			RegistryEnumsDirPath:      optEnumsDirPath,
			RegistryStructuresDirPath: optStructuresDirPath,
			RegistryEntitiesDirPath:   optEntitiesDirPath,
		},
		OutputPath:           outputDir,
		SourceCasing:         cfg.CasingSnake,
		TargetCasing:         cfg.CasingCamel,
		ZodSchemasImportPath: "@/schemas",
		TsTypesImportPath:    "@/types",
	}

	// Act
	err := compile.MorpheToBridge(config)
	s.NoError(err)

	// Verify User model: Email required, Nickname+Bio optional, Organization relationship optional
	userPath := filepath.Join(outputDir, "models", "user.ts")
	s.FileExists(userPath)
	userContent, readErr := os.ReadFile(userPath)
	s.NoError(readErr)
	userStr := string(userContent)

	// Required fields: direct assignment
	s.Contains(userStr, "    email: source.email,\n")
	s.Contains(userStr, "    id: source.id,\n")

	// Optional direct fields: conditional spread
	s.Contains(userStr, "...(source.bio !== undefined && {")
	s.Contains(userStr, "bio: source.bio,")
	s.Contains(userStr, "...(source.nickname !== undefined && {")
	s.Contains(userStr, "nickname: source.nickname,")

	// Relationship fields: always optional, with converter calls
	s.Contains(userStr, "import { convertOrganization } from './organization';")
	s.Contains(userStr, "...(source.organization_id !== undefined && {")
	s.Contains(userStr, "organizationID: source.organization_id,")
	s.Contains(userStr, "...(source.organization !== undefined && {")
	s.Contains(userStr, "organization: convertOrganization(source.organization),")
}
