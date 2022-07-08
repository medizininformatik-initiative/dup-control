package main

import (
	"fmt"
	"git.smith.care/smith/uc-phep/dupctl/cmd"
	"os"
)

func main() {
	rootCmd := cmd.NewRootCmd().Command()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
