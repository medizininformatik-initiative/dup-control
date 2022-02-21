package cli

import (
	"fmt"
	. "github.com/spf13/cobra"
)

// ExecutionError will return an Err and silence the cobra command usage() display.
func ExecutionError(cmd *Command, msg string, err error) error {
	cmd.SilenceUsage = true
	return fmt.Errorf(msg, err)
}
