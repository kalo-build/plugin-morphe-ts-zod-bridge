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
