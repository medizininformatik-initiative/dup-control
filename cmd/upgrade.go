package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
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
	Run: func(cmd *cobra.Command, args []string) {
		if err := updater.Upgrade(); err != nil {
			log.Infof("Error updating polarctl, %s", err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}
