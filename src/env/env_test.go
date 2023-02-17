package env

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type EnvSuite struct {
	suite.Suite
}

func TestEnvSuite(t *testing.T) {
	suite.Run(t, new(EnvSuite))
}

func (suite *EnvSuite) TestGetEnv() {
	suite.T().Setenv("TEST_ENV", "value-of-test-env")
	environmentVar := GetEnv("TEST_ENV", "")
	suite.Equal("value-of-test-env", environmentVar)
}

func (suite *EnvSuite) TestGetEnvWithDefaultValue() {
	environmentVar := GetEnv("TEST_ENV", "default-value")
	suite.Equal("default-value", environmentVar)
}
