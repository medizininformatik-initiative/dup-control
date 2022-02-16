package run

import (
	"fmt"
	coll "git.smith.care/smith/uc-phep/polar/polarctl/util"
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
	dev         bool
	env         map[string]string
}

var analyzeOpts = AnalyzeOpts{}

func createAnalyseOpts(analyzeOpts AnalyzeOpts) (container.PullOpts, container.RunOpts) {
	pullOpts := container.PullOpts{
		Image: fmt.Sprintf("%s-analysis", analyzeOpts.workpackage),
		Tag:   analyzeOpts.version,
	}
	runOpts := container.RunOpts{
		Env: coll.JoinEntries(analyzeOpts.env, "="),
		Mounts: []docker.HostMount{
			container.LocalMount("outputLocal", true),
			container.LocalMount("outputGlobal", true),
		},
	}
	if runtime.GOOS != "windows" {
		runOpts.User = fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid())
	}
	if analyzeOpts.dev {
		pullOpts.Image = "base"
		pullOpts.Tag = "latest"
		runOpts.Mounts = append(runOpts.Mounts,
			container.LocalMount("main.R", true),
			container.LocalMount("scripts", true),
			container.LocalMount("assets", true),
		)
	}
	return pullOpts, runOpts
}

var analyzeCommand = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze bundles retrieved from FHIR server",
	Long:  "You can analyze bundles that have formerly been retrieved from the FHIR server for a specific POLAR workpackage",
	PreRun: func(cmd *cobra.Command, args []string) {
		analyzeOpts.version = viper.GetString("analyze.version")
		analyzeOpts.env = viper.GetStringMapString("analyze.env")
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		pullOpts, runOpts := createAnalyseOpts(analyzeOpts)
		if !viper.GetBool("offline") {
			if err := containerRuntime.Pull(pullOpts); err != nil {
				return err
			}
		} else {
			log.Infof("Skip image pull due to --offline mode")
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

	analyzeCommand.PersistentFlags().BoolVar(&analyzeOpts.dev, "dev", false, "Mounts main.R, scripts/ and assets/ from current working directory for local development.")

	analyzeCommand.PersistentFlags().StringToStringP("env", "e", map[string]string{}, "Accepts key-value pairs in the form of key=value and passes them unchanged to the running scripts")
	_ = viper.BindPFlag("analyze.env", analyzeCommand.PersistentFlags().Lookup("env"))
}
