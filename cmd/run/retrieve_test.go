package run

import (
	"github.com/asaskevich/govalidator"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var rLog = logging.MustGetLogger("test")

func Test_createRetrieveOpts_GetsEnvInUppercase(t *testing.T) {
	viper.SetConfigFile("config.toml")
	err := viper.ReadConfig(strings.NewReader("retrieve.fhirServerEndpoint = \"https://host/\"\n retrieve.env = {\"FOO\" = \"BAR\"}"))
	if err != nil {
		rLog.Fatalf("Error Reading Config", err)
	}

	c := NewRetrieveCommand(rLog, nil)
	_ = c.Command().PreRunE(nil, nil)

	_, runOpts := c.createRetrieveOpts(c.retrieveOpts)
	for _, v := range runOpts.Env {
		rLog.Info(v)
		envKey := strings.Split(v, "=")[0]
		assert.True(t, govalidator.IsUpperCase(envKey), "Env var key should be uppercase")
	}
}
