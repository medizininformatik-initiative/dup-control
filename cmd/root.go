package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var cfgFile string
var vip = viper.New()

var workpackage string
var site string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "polarctl",
	Short:   "Control POLAR",
	Long:    `polarctl....`,
	Version: "0.1",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file")

	rootCmd.PersistentFlags().StringVar(&workpackage, "wp", "", "Workpackage to execute (e.g. 'wp-1-1-pilot').")
	_ = rootCmd.MarkPersistentFlagRequired("wp")

	rootCmd.PersistentFlags().StringVar(&site, "site", "latest", "Determines which image to use, as images are (not necessarily) hand-tailored for different dic sites. (e.g. 'dic-giessen', 'dic-leipzig', 'dic-muenchen').")
	_ = vip.BindPFlag("site", rootCmd.PersistentFlags().Lookup("site"))
}

func initConfig() {
	if cfgFile != "" {
		vip.SetConfigFile(cfgFile)
		if err := vip.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", vip.ConfigFileUsed())
		}
	}
}
