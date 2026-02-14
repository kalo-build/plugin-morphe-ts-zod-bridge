package cfg_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kalo-build/plugin-morphe-ts-zod-bridge/pkg/compile/cfg"
)

type CasingTestSuite struct {
	suite.Suite
}

func TestCasingTestSuite(t *testing.T) {
	suite.Run(t, new(CasingTestSuite))
}

func (s *CasingTestSuite) TestIsValid_ValidValues() {
	s.True(cfg.CasingCamel.IsValid())
	s.True(cfg.CasingSnake.IsValid())
	s.True(cfg.CasingPascal.IsValid())
}

func (s *CasingTestSuite) TestIsValid_InvalidValue() {
	s.False(cfg.Casing("kebab").IsValid())
	s.False(cfg.Casing("UPPER").IsValid())
	s.False(cfg.Casing("").IsValid())
}

func (s *CasingTestSuite) TestApply_CamelCase() {
	casing := cfg.CasingCamel
	s.Equal("firstName", casing.Apply("FirstName"))
	s.Equal("id", casing.Apply("ID"))
	s.Equal("taxID", casing.Apply("TaxID"))
}

func (s *CasingTestSuite) TestApply_SnakeCase() {
	casing := cfg.CasingSnake
	s.Equal("first_name", casing.Apply("FirstName"))
	s.Equal("id", casing.Apply("ID"))
	s.Equal("tax_id", casing.Apply("TaxID"))
}

func (s *CasingTestSuite) TestApply_PascalCase() {
	casing := cfg.CasingPascal
	s.Equal("FirstName", casing.Apply("FirstName"))
	s.Equal("ID", casing.Apply("ID"))
}

func (s *CasingTestSuite) TestErrInvalidCasing() {
	err := cfg.ErrInvalidCasing(cfg.Casing("foo"))
	s.Error(err)
	s.Contains(err.Error(), "foo")
	s.Contains(err.Error(), "invalid casing value")
}
