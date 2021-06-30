package cmd

import (
	"fmt"
	"git.smith.care/smith/uc-phep/polar/polarctl/util"
	"github.com/op/go-logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var log = logging.MustGetLogger("cmd")
var containerRuntime *container.Runtime
var cfgFile string
var workpackage string
var site string

var rootCmd = &cobra.Command{
	Use:     "polarctl",
	Short:   "Control POLAR",
	Long:    `polarctl....`,
	Version: "0.1",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		site = viper.GetString("site")

		if cr, err := container.NewRuntime(viper.GetString("registryUser"), viper.GetString("registryPass")); err != nil {
			return err
		} else {
			containerRuntime = cr
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

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config.toml", "Config file")

	rootCmd.PersistentFlags().StringVar(&workpackage, "wp", "", "Workpackage to execute (e.g. 'wp-1-1-pilot').")
	_ = rootCmd.MarkPersistentFlagRequired("wp")

	rootCmd.PersistentFlags().String("site", "latest", "Determines which image to use, as images are (not necessarily) hand-tailored for different dic sites. (e.g. 'dic-giessen', 'dic-leipzig', 'dic-muenchen').")
	_ = viper.BindPFlag("site", rootCmd.PersistentFlags().Lookup("site"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err == nil {
			log.Debugf("Using config file: %s", viper.ConfigFileUsed())
		}
	}
}
