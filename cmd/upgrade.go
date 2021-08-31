package cmd

import (
	"fmt"
	"git.smith.care/smith/uc-phep/polar/polarctl/util/upgrade"
	"github.com/spf13/cobra"
	"runtime"
)

func sprintReleaseKey(os string, arch string) string {
	if os != "windows" {
		return fmt.Sprintf("polarctl-%s-%s", os, arch)
	} else {
		return fmt.Sprintf("polarctl-%s-%s.exe", os, arch)
	}
}

func checkForUpdates() *upgrade.Updater {
	updater, err := upgrade.NewUpdater(baseURL, sprintReleaseKey(runtime.GOOS, runtime.GOARCH), "VERSION", Version)
	if err != nil {
		log.Warningf("Error checking for polarctl updates, %s", err.Error())
	} else {
		available, remoteVersion := updater.IsNewerVersionAvailable()
		if available {
			log.Infof("polarctl version %s available, use `polarctl upgrade` to download and replace your current version", remoteVersion)
		}
	}
	return updater
}

var upgradeCmd = &cobra.Command{
	Use:              "upgrade",
	Short:            "Upgrade polarctl",
	Long:             "You can upgrade your polarctl installation with the most recent version",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := updater.Upgrade(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}
