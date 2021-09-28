package cmd

import (
	"fmt"
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
	log              = logging.MustGetLogger("cmd")
	containerRuntime *container.Runtime
	cfgFile          string
	updater          *upgrade.Updater
	Version          = "v999.99"
)

var rootCmd = &cobra.Command{
	Use:     "polarctl",
	Short:   "Control POLAR",
	Long:    `polarctl....`,
	Version: Version,
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
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initUpdater() *upgrade.Updater {
	updater, err := upgrade.NewUpdater(baseURL, sprintReleaseKey(runtime.GOOS, runtime.GOARCH), "VERSION", Version)
	if err != nil {
		log.Fatalf("Error creating polarctl updater, %s", err.Error())
	}
	return updater
}

func checkForUpdates() {
	if !viper.GetBool("disableUpdateCheck") {
		available, remoteVersion := updater.IsNewerVersionAvailable()
		if available {
			log.Infof("polarctl version %s available, use `polarctl upgrade` to download and replace your current version", remoteVersion)
		}
	} else {
		log.Debugf("Upgrade checks disabled")
	}
}

func init() {
	updater = initUpdater()
	cobra.OnInitialize(initConfig, checkForUpdates)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config.toml", "Config file")
	rootCmd.PersistentFlags().Bool("disable-update-check", false, "Disable upgrade check on startup")
	_ = viper.BindPFlag("disableUpdateCheck", rootCmd.PersistentFlags().Lookup("disable-update-check"))
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
