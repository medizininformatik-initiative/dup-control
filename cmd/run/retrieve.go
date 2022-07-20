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
	"path/filepath"
	"runtime"
	"strings"
)

type retrieveOpts struct {
	dup                string //req
	version            string //req
	fhirServerEndpoint string
	fhirServerUser     string
	fhirServerPass     string
	fhirServerCACert   string
	fhirServerToken    string
	dev                bool
	env                map[string]string
}

type retrieveCommand struct {
	log *logging.Logger
	crp container.RuntimeProvider

	retrieveOpts retrieveOpts
}

func NewRetrieveCommand(log *logging.Logger, crp container.RuntimeProvider) *retrieveCommand {
	return &retrieveCommand{log: log, crp: crp, retrieveOpts: retrieveOpts{}}
}

func (c *retrieveCommand) createRetrieveOpts(retrieveOpts retrieveOpts) (container.PullOpts, container.RunOpts) {
	pullOpts := container.PullOpts{
		Image: retrieveOpts.dup,
		Tag:   retrieveOpts.version,
	}
	runOpts := container.RunOpts{
		Env: append(coll.JoinEntries(retrieveOpts.env, "="),
			fmt.Sprintf("FHIR_SERVER_ENDPOINT=%s", retrieveOpts.fhirServerEndpoint)),
		Mounts: []docker.HostMount{
			container.LocalMount("outputLocal", true),
			container.LocalMount("outputGlobal", true),
		},
	}
	if runtime.GOOS != "windows" {
		runOpts.User = fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid())
		runOpts.Env = append(runOpts.Env, "TERM=xterm-256color")
	}
	if retrieveOpts.fhirServerUser != "" && retrieveOpts.fhirServerPass != "" {
		runOpts.Env = append(runOpts.Env,
			fmt.Sprintf("FHIR_SERVER_USER=%s", retrieveOpts.fhirServerUser),
			fmt.Sprintf("FHIR_SERVER_PASS=%s", retrieveOpts.fhirServerPass))
	}
	if retrieveOpts.fhirServerCACert != "" {
		if abs, err := filepath.Abs(retrieveOpts.fhirServerCACert); err == nil {
			runOpts.Mounts = append(runOpts.Mounts,
				docker.HostMount{
					Source:   abs,
					Target:   "/etc/ssl/certs/ca-certificates.crt",
					Type:     "bind",
					ReadOnly: true})
		} else {
			c.log.Errorf("Skipping Certificate Injection: error converting path %s", retrieveOpts.fhirServerCACert)
		}
	}
	if retrieveOpts.fhirServerToken != "" {
		runOpts.Env = append(runOpts.Env,
			fmt.Sprintf("FHIR_SERVER_TOKEN=%s", retrieveOpts.fhirServerToken))
	}
	return pullOpts, runOpts
}

func (c *retrieveCommand) Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "retrieve",
		Short: "Retrieve bundles from FHIR server",
		Long:  "You can retrieve bundles from the FHIR server for a specific dup",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if viper.GetString("retrieve.fhirServerEndpoint") == "" {
				return fmt.Errorf("retrieve.fhirServerEndpoint not set")
			} else {
				c.retrieveOpts.fhirServerEndpoint = viper.GetString("retrieve.fhirServerEndpoint")
			}
			c.retrieveOpts.version = viper.GetString("retrieve.version")
			c.retrieveOpts.fhirServerUser = viper.GetString("retrieve.fhirServerUser")
			c.retrieveOpts.fhirServerPass = viper.GetString("retrieve.fhirServerPass")
			c.retrieveOpts.fhirServerCACert = viper.GetString("retrieve.fhirServerCACert")
			c.retrieveOpts.fhirServerToken = viper.GetString("retrieve.fhirServerToken")
			c.retrieveOpts.env = coll.TransformKeys(viper.GetStringMapString("retrieve.env"), strings.ToUpper)
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			containerRuntime, err := c.crp.CreateRuntime()
			if err != nil {
				return ExecutionError(cmd, "Unable to create ContainerRuntime, %w", err)
			}

			pullOpts, runOpts := c.createRetrieveOpts(c.retrieveOpts)
			if !viper.GetBool("offline") {
				if err := containerRuntime.Pull(pullOpts); err != nil {
					return ExecutionError(cmd, "Error pulling  retrieval image, %w", err)
				}
			} else {
				c.log.Infof("Skip image pull due to --offline mode")
			}

			if err := containerRuntime.Run("retrieval", pullOpts, runOpts); err != nil {
				return ExecutionError(cmd, "Error running retrieval container, %w", err)
			}

			return nil
		},
	}

	command.PersistentFlags().StringVar(&c.retrieveOpts.dup, "dup", "", "Image to execute (e.g. 'vhf').")
	_ = command.MarkPersistentFlagRequired("dup")

	command.PersistentFlags().String("version", "latest", "Determines which image to use, as images can be versioned or hand-tailored for different dic sites. (e.g. '0.1', 'dic-giessen', 'dic-leipzig', 'dic-muenchen').")
	_ = viper.BindPFlag("retrieve.version", command.PersistentFlags().Lookup("version"))

	command.PersistentFlags().String("fhir-server-endpoint", "", "the base URL of the FHIR server to use")
	_ = viper.BindPFlag("retrieve.fhirServerEndpoint", command.PersistentFlags().Lookup("fhir-server-endpoint"))

	command.PersistentFlags().String("fhir-server-user", "", "Username for basic auth protected communication with FHIR server")
	_ = viper.BindPFlag("retrieve.fhirServerUser", command.PersistentFlags().Lookup("fhir-server-user"))

	command.PersistentFlags().String("fhir-server-pass", "", "Password for basic auth protected communication with FHIR server")
	_ = viper.BindPFlag("retrieve.fhirServerPass", command.PersistentFlags().Lookup("fhir-server-pass"))

	command.PersistentFlags().String("fhir-server-token", "", "Token for token based auth protected communication with FHIR Server")
	_ = viper.BindPFlag("retrieve.fhirServerToken", command.PersistentFlags().Lookup("fhir-server-token"))

	command.PersistentFlags().StringToStringP("env", "e", map[string]string{}, "Accepts key-value pairs in the form of key=value and passes them unchanged to the running scripts")
	_ = viper.BindPFlag("retrieve.env", command.PersistentFlags().Lookup("env"))

	return command
}
