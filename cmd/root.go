package cmd

import (
	"fmt"
	"git.smith.care/smith/uc-phep/dupctl/cmd/run"
	"git.smith.care/smith/uc-phep/dupctl/cmd/util"
	. "git.smith.care/smith/uc-phep/dupctl/lib/cli"
	"git.smith.care/smith/uc-phep/dupctl/lib/container"
	"git.smith.care/smith/uc-phep/dupctl/lib/upgrade"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/op/go-logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"runtime"
)

var Version = "v999.99.99-dev"

type rootCmd struct {
	log            *logging.Logger
	version        string
	updaterBaseUrl string
	cfgFile        string
}

func NewRootCmd() *rootCmd {
	return &rootCmd{
		log:            logging.MustGetLogger("cmd"),
		version:        Version,
		updaterBaseUrl: "https://dupctl.s3.amazonaws.com",
	}
}

func (c *rootCmd) Command() *cobra.Command {
	updater := c.initUpdater()

	command := &cobra.Command{
		Use:   "dupctl",
		Short: "Execute Data Use Project containers",
		Long:  `dupctl....`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := c.initConfig(); err != nil {
				return ExecutionError(cmd, "error initializing config, %w", err)
			}
			c.checkForUpdates(updater)
			return nil
		},
		Version: c.version,
	}

	command.PersistentFlags().StringVar(&c.cfgFile, "config", "config.toml", "Config file")
	command.PersistentFlags().Bool("disable-update-check", false, "Disable upgrade check on startup")
	_ = viper.BindPFlag("disableUpdateCheck", command.PersistentFlags().Lookup("disable-update-check"))
	command.PersistentFlags().Bool("offline", false, "Assumes an air-gapped environment.")
	_ = viper.BindPFlag("offline", command.PersistentFlags().Lookup("offline"))

	crp := &containerRuntimeProvider{}
	command.AddCommand(run.NewAnalyzeCommand(c.log, crp).Command())
	command.AddCommand(run.NewRetrieveCommand(c.log, crp).Command())
	command.AddCommand(util.NewUpgradeCommand(c.log, updater).Command())

	return command
}

func (c *rootCmd) initConfig() error {
	if c.cfgFile != "" {
		viper.SetConfigFile(c.cfgFile)
		if err := viper.ReadInConfig(); err != nil {
			return fmt.Errorf("error reading config file %s, %w", viper.ConfigFileUsed(), err)
		}
		c.log.Debugf("Using config file: %s", viper.ConfigFileUsed())
	}
	return nil
}

func (c *rootCmd) initUpdater() upgrade.Updater {
	updater, err := upgrade.NewUpdater(c.updaterBaseUrl, runtime.GOOS, runtime.GOARCH, "VERSION", c.version)
	if err != nil {
		log.Fatalf("Error creating dupctl updater: %v", err)
	}
	return updater
}

func (c *rootCmd) checkForUpdates(updater upgrade.Updater) {
	if !viper.GetBool("disableUpdateCheck") && !viper.GetBool("offline") {
		available, remoteVersion := updater.IsNewerVersionAvailable()
		if available {
			c.log.Infof("dupctl version %s available, use `dupctl upgrade` to download and replace your current version", remoteVersion)
		}
	} else {
		c.log.Debugf("dupctl checks disabled")
	}
}

type containerRuntimeProvider struct {
}

func (p *containerRuntimeProvider) CreateRuntime() (*container.Runtime, error) {
	if viper.GetString("registryUser") == "" {
		return nil, fmt.Errorf("registryUser not configured. Please check if you are missing the 'config.toml' config file")
	} else if viper.GetString("registryPass") == "" {
		return nil, fmt.Errorf("registryPass not configured. Please check if you are missing the 'config.toml' config file")
	} else if cli, err := docker.NewClientFromEnv(); err != nil {
		return nil, fmt.Errorf("cannot instantiate docker client, %w", err)
	} else {
		return container.NewRuntime(cli, viper.GetString("registry"), viper.GetString("project"),
			viper.GetString("registryUser"), viper.GetString("registryPass")), nil
	}
}
