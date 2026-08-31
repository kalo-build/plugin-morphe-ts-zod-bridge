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

type EntityConverterTestSuite struct {
	suite.Suite

	TestDirPath       string
	ModelsDirPath     string
	EnumsDirPath      string
	StructuresDirPath string
	EntitiesDirPath   string
}

func TestEntityConverterTestSuite(t *testing.T) {
	suite.Run(t, new(EntityConverterTestSuite))
}

func (s *EntityConverterTestSuite) SetupTest() {
	s.TestDirPath = testutils.GetTestDirPath()
	s.ModelsDirPath = filepath.Join(s.TestDirPath, "registry", "minimal", "models")
	s.EnumsDirPath = filepath.Join(s.TestDirPath, "registry", "minimal", "enums")
	s.StructuresDirPath = filepath.Join(s.TestDirPath, "registry", "minimal", "structures")
	s.EntitiesDirPath = filepath.Join(s.TestDirPath, "registry", "minimal", "entities")
}

func (s *EntityConverterTestSuite) TestCompileAllEntityConverters_SnakeToCamel() {
	// Arrange
	outputDir := filepath.Join(s.TestDirPath, "working-entities-snake-camel")
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

	// Verify Person entity converter
	personPath := filepath.Join(outputDir, "entities", "person.ts")
	s.FileExists(personPath)
	personContent, readErr := os.ReadFile(personPath)
	s.NoError(readErr)
	personStr := string(personContent)

	s.Contains(personStr, "import type { Person } from '@/generated/types/entities/person';")
	s.Contains(personStr, "import { PersonSchema } from '@/generated/schemas/entities/person';")
	s.Contains(personStr, "export function convertPerson(")
	s.Contains(personStr, "email: source.email,")
	s.Contains(personStr, "id: source.id,")
	s.Contains(personStr, "lastName: source.last_name,")
	s.Contains(personStr, "nationality: source.nationality,")

	// Verify Company entity converter
	companyPath := filepath.Join(outputDir, "entities", "company.ts")
	s.FileExists(companyPath)
	companyContent, readErr := os.ReadFile(companyPath)
	s.NoError(readErr)
	companyStr := string(companyContent)

	s.Contains(companyStr, "export function convertCompany(")
	s.Contains(companyStr, "name: source.name,")
	s.Contains(companyStr, "taxID: source.tax_id,")
}

func (s *EntityConverterTestSuite) TestCompileAllEntityConverters_RelationshipFields() {
	// Arrange
	outputDir := filepath.Join(s.TestDirPath, "working-entities-relationships")
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

	// Person entity has ForOne Company -- produces optional FK + object fields
	personPath := filepath.Join(outputDir, "entities", "person.ts")
	personContent, readErr := os.ReadFile(personPath)
	s.NoError(readErr)
	personStr := string(personContent)

	// Required FK field (ForOne Company without optional): direct assignment
	s.Contains(personStr, "import { convertCompany } from './company';")
	s.Contains(personStr, "    companyID: source.company_id,\n")
	// Relation object: always optional
	s.Contains(personStr, "...(source.company !== undefined && {")
	s.Contains(personStr, "company: convertCompany(source.company),")

	// Company entity has HasMany Person -- IDs required, objects optional
	companyPath := filepath.Join(outputDir, "entities", "company.ts")
	companyContent, readErr := os.ReadFile(companyPath)
	s.NoError(readErr)
	companyStr := string(companyContent)

	s.Contains(companyStr, "import { convertPerson } from './person';")
	s.Contains(companyStr, "    personIDs: source.person_ids,\n")
	s.Contains(companyStr, "...(source.persons !== undefined && {")
	s.Contains(companyStr, "persons: source.persons.map(convertPerson),")
}

func (s *EntityConverterTestSuite) TestCompileAllEntityConverters_PolymorphicRelationships() {
	// Arrange — uses the polymorphic testdata registry
	polyModelsDirPath := filepath.Join(s.TestDirPath, "registry", "polymorphic", "models")
	polyEnumsDirPath := filepath.Join(s.TestDirPath, "registry", "polymorphic", "enums")
	polyStructuresDirPath := filepath.Join(s.TestDirPath, "registry", "polymorphic", "structures")
	polyEntitiesDirPath := filepath.Join(s.TestDirPath, "registry", "polymorphic", "entities")

	outputDir := filepath.Join(s.TestDirPath, "working-entities-poly")
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

	// Assert — Comment entity has ForOnePoly Commentable [Person, Company]
	commentPath := filepath.Join(outputDir, "entities", "comment.ts")
	s.FileExists(commentPath)
	commentContent, readErr := os.ReadFile(commentPath)
	s.NoError(readErr)
	commentStr := string(commentContent)

	// ForOnePoly produces discriminator fields (required unless relation has attributes: [optional])
	s.Contains(commentStr, "commentableID: source.commentable_id,")
	s.Contains(commentStr, "commentableType: source.commentable_type,")

	// Assert — Person entity has HasOnePoly Note (Aliased: Comment)
	personPath := filepath.Join(outputDir, "entities", "person.ts")
	s.FileExists(personPath)
	personContent, readErr := os.ReadFile(personPath)
	s.NoError(readErr)
	personStr := string(personContent)

	// HasOnePoly produces FK (required, direct assignment) + object reference (always optional)
	s.Contains(personStr, "import { convertComment } from './comment';")
	s.Contains(personStr, "    noteID: source.note_id,\n")
	s.Contains(personStr, "...(source.note !== undefined && {")
	s.Contains(personStr, "note: convertComment(source.note),")
}

func (s *EntityConverterTestSuite) TestCompileAllEntityConverters_OptionalFields() {
	// Arrange -- uses the optional testdata registry
	optModelsDirPath := filepath.Join(s.TestDirPath, "registry", "optional", "models")
	optEnumsDirPath := filepath.Join(s.TestDirPath, "registry", "optional", "enums")
	optStructuresDirPath := filepath.Join(s.TestDirPath, "registry", "optional", "structures")
	optEntitiesDirPath := filepath.Join(s.TestDirPath, "registry", "optional", "entities")

	outputDir := filepath.Join(s.TestDirPath, "working-entities-optional")
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

	// Verify User entity: Email+ID required, Nickname optional, Organization relationship optional
	userPath := filepath.Join(outputDir, "entities", "user.ts")
	s.FileExists(userPath)
	userContent, readErr := os.ReadFile(userPath)
	s.NoError(readErr)
	userStr := string(userContent)

	// Required fields
	s.Contains(userStr, "    email: source.email,\n")
	s.Contains(userStr, "    id: source.id,\n")

	// Optional direct field
	s.Contains(userStr, "...(source.nickname !== undefined && {")
	s.Contains(userStr, "nickname: source.nickname,")

	// Required FK field (Organization ForOne without optional): direct assignment
	s.Contains(userStr, "import { convertOrganization } from './organization';")
	s.Contains(userStr, "    organizationID: source.organization_id,\n")
	// Relation object: always optional
	s.Contains(userStr, "...(source.organization !== undefined && {")
	s.Contains(userStr, "organization: convertOrganization(source.organization),")
}

func (s *EntityConverterTestSuite) TestCompileAllEntityConverters_DifferentImportPaths() {
	// Arrange
	outputDir := filepath.Join(s.TestDirPath, "working-entities-custom-paths")
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
		ZodSchemasImportPath: "~/api/schemas",
		TsTypesImportPath:    "~/api/types",
	}

	// Act
	err := compile.MorpheToBridge(config)

	// Assert
	s.NoError(err)

	personPath := filepath.Join(outputDir, "entities", "person.ts")
	s.FileExists(personPath)
	personContent, readErr := os.ReadFile(personPath)
	s.NoError(readErr)
	personStr := string(personContent)

	s.Contains(personStr, "import type { Person } from '~/api/types/entities/person';")
	s.Contains(personStr, "import { PersonSchema } from '~/api/schemas/entities/person';")
}
