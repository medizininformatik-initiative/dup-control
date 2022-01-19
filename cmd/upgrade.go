package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

func sprintReleaseKey(os string, arch string) string {
	if os != "windows" {
		return fmt.Sprintf("polarctl-%s-%s", os, arch)
	} else {
		return fmt.Sprintf("polarctl-%s-%s.exe", os, arch)
	}
}

var upgradeCmd = &cobra.Command{
	Use:              "upgrade",
	Short:            "Upgrade polarctl",
	Long:             "You can upgrade your polarctl installation with the most recent version",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {},
	RunE: func(cmd *cobra.Command, args []string) error {
		if !viper.GetBool("offline") {
			if err := updater.Upgrade(); err != nil {
				log.Infof("Error updating polarctl: %v", err)
				os.Exit(1)
			}
		} else {
			return fmt.Errorf("cannot upgrade in --offline mode")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}
