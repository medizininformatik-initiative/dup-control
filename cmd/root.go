package cmd

import (
	"fmt"
	"git.smith.care/smith/uc-phep/polar/polarctl/util/container"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/op/go-logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var log = logging.MustGetLogger("cmd")
var containerRuntime *container.Runtime
var cfgFile string
var rootOpts = RootOpts{}

type RootOpts struct{}

var rootCmd = &cobra.Command{
	Use:   "polarctl",
	Short: "Control POLAR",
	Long:  `polarctl....`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cli, err := docker.NewClientFromEnv(); err != nil {
			return fmt.Errorf("cannot instantiate docker client, %w", err)
		} else {
			containerRuntime = container.NewRuntime(cli, viper.GetString("registryUser"), viper.GetString("registryPass"))
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	rootCmd.Version = version
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config.toml", "Config file")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err == nil {
			log.Debugf("Using config file: %s", viper.ConfigFileUsed())
		} else {
			log.Errorf("Error reading config file: %s", err.Error())
			os.Exit(3)
		}
	}
}
