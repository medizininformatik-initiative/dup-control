package main

import (
	"fmt"
	"git.smith.care/smith/uc-phep/polar/polarctl/cmd"
	"os"
)

func main() {
	rootCmd := cmd.NewRootCmd().Command()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
