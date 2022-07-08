package run

import (
	"fmt"
	. "git.smith.care/smith/uc-phep/dupctl/lib/cli"
	"git.smith.care/smith/uc-phep/dupctl/lib/coll"
	"git.smith.care/smith/uc-phep/dupctl/lib/container"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/op/go-logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"runtime"
)

type analyzeOpts struct {
	workpackage string //req
	version     string
	dev         bool
	env         map[string]string
}

type analyzeCommand struct {
	log *logging.Logger
	crp container.RuntimeProvider

	analyzeOpts analyzeOpts
}

func NewAnalyzeCommand(log *logging.Logger, crp container.RuntimeProvider) *analyzeCommand {
	return &analyzeCommand{log: log, crp: crp, analyzeOpts: analyzeOpts{}}
}

func (c *analyzeCommand) createAnalyseOpts(analyzeOpts analyzeOpts) (container.PullOpts, container.RunOpts) {
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
		runOpts.Env = append(runOpts.Env, "TERM=xterm-256color")
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

func (c *analyzeCommand) Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "analyze",
		Short: "Analyze bundles retrieved from FHIR server",
		Long:  "You can analyze bundles that have formerly been retrieved from the FHIR server for a specific dup",
		PreRun: func(cmd *cobra.Command, args []string) {
			c.analyzeOpts.version = viper.GetString("analyze.version")
			c.analyzeOpts.env = viper.GetStringMapString("analyze.env")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			containerRuntime, err := c.crp.CreateRuntime()
			if err != nil {
				return ExecutionError(cmd, "Unable to create ContainerRuntime, %w", err)
			}

			pullOpts, runOpts := c.createAnalyseOpts(c.analyzeOpts)
			if !viper.GetBool("offline") {
				if err := containerRuntime.Pull(pullOpts); err != nil {
					return ExecutionError(cmd, "Error pulling analysis image, %w", err)
				}
			} else {
				c.log.Infof("Skip image pull due to --offline mode")
			}

			if err := containerRuntime.Run("analysis", pullOpts, runOpts); err != nil {
				return ExecutionError(cmd, "Error running analysis container, %w", err)
			}

			return nil
		},
	}

	command.PersistentFlags().StringVar(&c.analyzeOpts.workpackage, "wp", "", "Image to execute (e.g. 'wp-1-1-pilot').")
	_ = command.MarkPersistentFlagRequired("wp")

	command.PersistentFlags().String("version", "latest", "Determines which image version to use.")
	_ = viper.BindPFlag("analyze.version", command.PersistentFlags().Lookup("version"))

	command.PersistentFlags().BoolVar(&c.analyzeOpts.dev, "dev", false, "Mounts main.R, scripts/ and assets/ from current working directory for local development.")

	command.PersistentFlags().StringToStringP("env", "e", map[string]string{}, "Accepts key-value pairs in the form of key=value and passes them unchanged to the running scripts")
	_ = viper.BindPFlag("analyze.env", command.PersistentFlags().Lookup("env"))

	return command
}
