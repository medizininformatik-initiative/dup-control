package cmd

import (
	"fmt"
	"git.smith.care/smith/uc-phep/polar/polarctl/util/container"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"runtime"
)

type AnalyzeOpts struct {
	workpackage string //req
	version     string
}

var analyzeOpts = AnalyzeOpts{}

func createAnalyseOpts(analyzeOpts AnalyzeOpts) (container.PullOpts, container.RunOpts) {
	pullOpts := container.PullOpts{
		Image: fmt.Sprintf("%s-analysis", analyzeOpts.workpackage),
		Tag:   analyzeOpts.version,
	}
	runOpts := container.RunOpts{
		Env: []string{},
		Mounts: []docker.HostMount{
			container.LocalMount("outputLocal", true),
			container.LocalMount("outputGlobal", true),
		},
	}
	if runtime.GOOS != "windows" {
		runOpts.User = fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid())
	}
	return pullOpts, runOpts
}

var analyzeCommand = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze bundles retrieved from FHIR server",
	Long:  "You can analyze bundles that have formerly been retrieved from the FHIR server for a specific POLAR workpackage",
	PreRun: func(cmd *cobra.Command, args []string) {
		analyzeOpts.version = viper.GetString("analyze.version")
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		pullOpts, runOpts := createAnalyseOpts(analyzeOpts)
		if err := containerRuntime.Pull(pullOpts); err != nil {
			return err
		}

		if err := containerRuntime.Run("analysis", pullOpts, runOpts); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(analyzeCommand)

	analyzeCommand.PersistentFlags().StringVar(&analyzeOpts.workpackage, "wp", "", "Image to execute (e.g. 'wp-1-1-pilot').")
	_ = analyzeCommand.MarkPersistentFlagRequired("wp")

	analyzeCommand.PersistentFlags().String("version", "latest", "Determines which image version to use.")
	_ = viper.BindPFlag("analyze.version", analyzeCommand.PersistentFlags().Lookup("version"))
}
