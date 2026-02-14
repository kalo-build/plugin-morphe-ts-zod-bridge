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
