package util

import (
	"fmt"
	"git.smith.care/smith/uc-phep/polar/polarctl/util/upgrade"
	"github.com/op/go-logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

type upgradeCommand struct {
	log     *logging.Logger
	updater upgrade.Updater
}

func NewUpgradeCommand(log *logging.Logger, updater upgrade.Updater) *upgradeCommand {
	return &upgradeCommand{log: log, updater: updater}
}

func (c *upgradeCommand) Command() *cobra.Command {
	return &cobra.Command{
		Use:              "upgrade",
		Short:            "Upgrade polarctl",
		Long:             "You can upgrade your polarctl installation with the most recent version",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {},
		RunE: func(cmd *cobra.Command, args []string) error {
			if !viper.GetBool("offline") {
				if err := c.updater.Upgrade(); err != nil {
					c.log.Infof("Error updating polarctl: %v", err)
					os.Exit(1)
				}
			} else {
				return fmt.Errorf("cannot upgrade in --offline mode")
			}
			return nil
		},
	}
}
