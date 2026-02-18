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

type StructureConverterTestSuite struct {
	suite.Suite

	TestDirPath       string
	ModelsDirPath     string
	EnumsDirPath      string
	StructuresDirPath string
	EntitiesDirPath   string
}

func TestStructureConverterTestSuite(t *testing.T) {
	suite.Run(t, new(StructureConverterTestSuite))
}

func (s *StructureConverterTestSuite) SetupTest() {
	s.TestDirPath = testutils.GetTestDirPath()
	s.ModelsDirPath = filepath.Join(s.TestDirPath, "registry", "minimal", "models")
	s.EnumsDirPath = filepath.Join(s.TestDirPath, "registry", "minimal", "enums")
	s.StructuresDirPath = filepath.Join(s.TestDirPath, "registry", "minimal", "structures")
	s.EntitiesDirPath = filepath.Join(s.TestDirPath, "registry", "minimal", "entities")
}

func (s *StructureConverterTestSuite) TestCompileAllStructureConverters_SnakeToCamel() {
	// Arrange
	outputDir := filepath.Join(s.TestDirPath, "working-structures-snake-camel")
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

	// Verify Address structure converter
	addressPath := filepath.Join(outputDir, "structures", "address.ts")
	s.FileExists(addressPath)
	addressContent, readErr := os.ReadFile(addressPath)
	s.NoError(readErr)
	addressStr := string(addressContent)

	s.Contains(addressStr, "import type { Address } from '@/generated/types/structures/address';")
	s.Contains(addressStr, "import { AddressSchema } from '@/generated/schemas/structures/address';")
	s.Contains(addressStr, "export function convertAddress(")
	s.Contains(addressStr, "source: z.infer<typeof AddressSchema>")
	s.Contains(addressStr, "): Address {")
	s.Contains(addressStr, "city: source.city,")
	s.Contains(addressStr, "houseNr: source.house_nr,")
	s.Contains(addressStr, "street: source.street,")
	s.Contains(addressStr, "zipCode: source.zip_code,")
}

func (s *StructureConverterTestSuite) TestCompileAllStructureConverters_PascalToSnake() {
	// Arrange
	outputDir := filepath.Join(s.TestDirPath, "working-structures-pascal-snake")
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
		SourceCasing:         cfg.CasingPascal,
		TargetCasing:         cfg.CasingSnake,
		ZodSchemasImportPath: "@/schemas",
		TsTypesImportPath:    "@/types",
	}

	// Act
	err := compile.MorpheToBridge(config)

	// Assert
	s.NoError(err)

	addressPath := filepath.Join(outputDir, "structures", "address.ts")
	s.FileExists(addressPath)
	addressContent, readErr := os.ReadFile(addressPath)
	s.NoError(readErr)
	addressStr := string(addressContent)

	// Source is PascalCase, target is snake_case
	s.Contains(addressStr, "house_nr: source.HouseNr,")
	s.Contains(addressStr, "zip_code: source.ZipCode,")
}

func (s *StructureConverterTestSuite) TestCompileAllStructureConverters_OptionalFields() {
	// Arrange -- uses the optional testdata registry
	optModelsDirPath := filepath.Join(s.TestDirPath, "registry", "optional", "models")
	optEnumsDirPath := filepath.Join(s.TestDirPath, "registry", "optional", "enums")
	optStructuresDirPath := filepath.Join(s.TestDirPath, "registry", "optional", "structures")
	optEntitiesDirPath := filepath.Join(s.TestDirPath, "registry", "optional", "entities")

	outputDir := filepath.Join(s.TestDirPath, "working-structures-optional")
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

	// Verify ProfileResponse: Email+DisplayName required, AvatarURL optional
	profilePath := filepath.Join(outputDir, "structures", "profile-response.ts")
	s.FileExists(profilePath)
	profileContent, readErr := os.ReadFile(profilePath)
	s.NoError(readErr)
	profileStr := string(profileContent)

	// Required fields: direct assignment
	s.Contains(profileStr, "    displayName: source.display_name,\n")
	s.Contains(profileStr, "    email: source.email,\n")

	// Optional field: conditional spread
	s.Contains(profileStr, "...(source.avatar_url !== undefined && {")
	s.Contains(profileStr, "avatarURL: source.avatar_url,")
}
