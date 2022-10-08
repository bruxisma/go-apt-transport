package transport

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConfigurationSuite struct {
	suite.Suite
}

func (suite *ConfigurationSuite) TestConfigurationUnmarshalFields() {
	fields := Fields{
		"Config-Item": []string{
			"APT::Install-Recommends=false",
			"APT::Get::Assume-Yes=true",
		},
	}
	configuration := Configuration{}
	err := UnmarshalFields(fields, configuration)
	suite.NoError(err)
	suite.Contains(configuration, "APT::Install-Recommends")
	suite.Contains(configuration.Section("APT"), "Get::Assume-Yes")
}

func TestHandshakeMessages(test *testing.T) {
	suite.Run(test, new(ConfigurationSuite))
}
