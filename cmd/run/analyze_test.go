package run

import (
	"github.com/asaskevich/govalidator"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var aLog = logging.MustGetLogger("test")

func Test_createAnalyzeOpts_GetsEnvInUppercase(t *testing.T) {
	viper.SetConfigFile("config.toml")
	err := viper.ReadConfig(strings.NewReader("analyze.env = {\"FOO\" = \"BAR\"}"))
	if err != nil {
		aLog.Fatalf("Error Reading Config", err)
	}

	c := NewAnalyzeCommand(aLog, nil)
	c.Command().PreRun(nil, nil)

	_, runOpts := c.createAnalyseOpts(c.analyzeOpts)
	for _, v := range runOpts.Env {
		aLog.Info(v)
		envKey := strings.Split(v, "=")[0]
		assert.True(t, govalidator.IsUpperCase(envKey), "Env var key should be uppercase")
	}
}
