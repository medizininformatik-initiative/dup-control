package cmd

import (
	"fmt"
	"git.smith.care/smith/uc-phep/polar/polarctl/cmd/run"
	"git.smith.care/smith/uc-phep/polar/polarctl/cmd/util"
	"git.smith.care/smith/uc-phep/polar/polarctl/util/container"
	"git.smith.care/smith/uc-phep/polar/polarctl/util/upgrade"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/op/go-logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"runtime"
)

const baseURL = "https://polarctl.s3.amazonaws.com"

var (
	log     = logging.MustGetLogger("cmd")
	cfgFile string
	Version = "v999.99.99-dev"
)

var rootCmd = &cobra.Command{
	Use:     "polarctl",
	Short:   "Control POLAR",
	Long:    `polarctl....`,
	Version: Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initUpdater() *upgrade.Updater {
	updater, err := upgrade.NewUpdater(baseURL, runtime.GOOS, runtime.GOARCH, "VERSION", Version)
	if err != nil {
		log.Fatalf("Error creating polarctl updater: %v", err)
	}
	return updater
}

func checkForUpdates(updater *upgrade.Updater) {
	if !viper.GetBool("disableUpdateCheck") && !viper.GetBool("offline") {
		available, remoteVersion := updater.IsNewerVersionAvailable()
		if available {
			log.Infof("polarctl version %s available, use `polarctl upgrade` to download and replace your current version", remoteVersion)
		}
	} else {
		log.Debugf("Upgrade checks disabled")
	}
}

func init() {
	updater := initUpdater()
	checkForUpdates(updater)
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config.toml", "Config file")
	rootCmd.PersistentFlags().Bool("disable-update-check", false, "Disable upgrade check on startup")
	_ = viper.BindPFlag("disableUpdateCheck", rootCmd.PersistentFlags().Lookup("disable-update-check"))
	rootCmd.PersistentFlags().Bool("offline", false, "Assumes an air-gapped environment.")
	_ = viper.BindPFlag("offline", rootCmd.PersistentFlags().Lookup("offline"))

	provider := &runtimeProvider{}
	rootCmd.AddCommand(run.AnalyzeCommand(log, provider))
	rootCmd.AddCommand(run.RetrieveCommand(log, provider))
	rootCmd.AddCommand(util.UpgradeCommand(log, updater))
}

type runtimeProvider struct {
}

func (p *runtimeProvider) CreateRuntime() (*container.Runtime, error) {
	if viper.GetString("registryUser") == "" {
		return nil, fmt.Errorf("registryUser not configured. Please check if you are missing the 'config.toml' config file")
	} else if viper.GetString("registryPass") == "" {
		return nil, fmt.Errorf("registryPass not configured. Please check if you are missing the 'config.toml' config file")
	} else if cli, err := docker.NewClientFromEnv(); err != nil {
		return nil, fmt.Errorf("cannot instantiate docker client, %w", err)
	} else {
		return container.NewRuntime(cli, viper.GetString("registryUser"), viper.GetString("registryPass")), nil
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err == nil {
			log.Debugf("Using config file: %s", viper.ConfigFileUsed())
		}
	}
}
