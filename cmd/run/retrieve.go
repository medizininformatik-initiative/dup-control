package run

import (
	"fmt"
	coll "git.smith.care/smith/uc-phep/polar/polarctl/util"
	"git.smith.care/smith/uc-phep/polar/polarctl/util/container"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/op/go-logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"runtime"
)

type RetrieveOpts struct {
	workpackage        string //req
	version            string //req
	fhirServerEndpoint string
	fhirServerUser     string
	fhirServerPass     string
	fhirServerCACert   string
	fhirServerToken    string
	dev                bool
	env                map[string]string
}

func RetrieveCommand(log *logging.Logger, crp container.RuntimeProvider) *cobra.Command {
	var retrieveOpts = RetrieveOpts{}

	createRetrieveOpts := func(retrieveOpts RetrieveOpts) (container.PullOpts, container.RunOpts) {
		pullOpts := container.PullOpts{
			Image: retrieveOpts.workpackage,
			Tag:   retrieveOpts.version,
		}
		runOpts := container.RunOpts{
			Env: append(coll.JoinEntries(retrieveOpts.env, "="),
				fmt.Sprintf("FHIR_ENDPOINT=%s", retrieveOpts.fhirServerEndpoint)),
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
				fmt.Sprintf("FHIR_USERNAME=%s", retrieveOpts.fhirServerUser),
				fmt.Sprintf("FHIR_PASSWORD=%s", retrieveOpts.fhirServerPass))
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
				log.Errorf("Skipping Certificate Injection: error converting path %s", retrieveOpts.fhirServerCACert)
			}
		}
		if retrieveOpts.fhirServerToken != "" {
			runOpts.Env = append(runOpts.Env,
				fmt.Sprintf("FHIR_TOKEN=%s", retrieveOpts.fhirServerToken))
		}
		if retrieveOpts.dev {
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

	command := &cobra.Command{
		Use:   "retrieve",
		Short: "Retrieve bundles from FHIR server",
		Long:  "You can retrieve bundles from the FHIR server for a specific POLAR workpackage",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if viper.GetString("retrieve.fhirServerEndpoint") == "" {
				return fmt.Errorf("retrieve.fhirServerEndpoint not set")
			} else {
				retrieveOpts.fhirServerEndpoint = viper.GetString("retrieve.fhirServerEndpoint")
			}
			if site := viper.GetString("retrieve.site"); site != "latest" {
				log.Warningf("--site flag / retrieve.site option is deprecated! Use --version flag / retrieve.version option instead!")
				retrieveOpts.version = site
			} else {
				retrieveOpts.version = viper.GetString("retrieve.version")
			}
			retrieveOpts.fhirServerUser = viper.GetString("retrieve.fhirServerUser")
			retrieveOpts.fhirServerPass = viper.GetString("retrieve.fhirServerPass")
			retrieveOpts.fhirServerCACert = viper.GetString("retrieve.fhirServerCACert")
			retrieveOpts.fhirServerToken = viper.GetString("retrieve.fhirServerToken")
			retrieveOpts.env = viper.GetStringMapString("retrieve.env")
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			containerRuntime, err := crp.CreateRuntime()
			if err != nil {
				log.Fatalf("Unable to create ContainerRuntime, %w", err)
			}

			pullOpts, runOpts := createRetrieveOpts(retrieveOpts)
			if !viper.GetBool("offline") {
				if err := containerRuntime.Pull(pullOpts); err != nil {
					return err
				}
			} else {
				log.Infof("Skip image pull due to --offline mode")
			}

			if err := containerRuntime.Run("retrieval", pullOpts, runOpts); err != nil {
				return err
			}

			return nil
		},
	}

	command.PersistentFlags().StringVar(&retrieveOpts.workpackage, "wp", "", "Image to execute (e.g. 'wp-1-1-pilot').")
	_ = command.MarkPersistentFlagRequired("wp")

	command.PersistentFlags().String("version", "latest", "Determines which image to use, as images can be versioned or hand-tailored for different dic sites. (e.g. '0.1', 'dic-giessen', 'dic-leipzig', 'dic-muenchen').")
	_ = viper.BindPFlag("retrieve.version", command.PersistentFlags().Lookup("version"))
	command.PersistentFlags().String("site", "latest", "Determines which image to use, as images can be hand-tailored for different dic sites. (e.g. 'dic-giessen', 'dic-leipzig', 'dic-muenchen'). DEPRECATED! Use --version instead!")
	_ = viper.BindPFlag("retrieve.site", command.PersistentFlags().Lookup("site"))

	command.PersistentFlags().String("fhir-server-endpoint", "", "the base URL of the FHIR server to use")
	_ = viper.BindPFlag("retrieve.fhirServerEndpoint", command.PersistentFlags().Lookup("fhir-server-endpoint"))

	command.PersistentFlags().String("fhir-server-user", "", "Username for basic auth protected communication with FHIR server")
	_ = viper.BindPFlag("retrieve.fhirServerUser", command.PersistentFlags().Lookup("fhir-server-user"))

	command.PersistentFlags().String("fhir-server-pass", "", "Password for basic auth protected communication with FHIR server")
	_ = viper.BindPFlag("retrieve.fhirServerPass", command.PersistentFlags().Lookup("fhir-server-pass"))

	command.PersistentFlags().String("fhir-server-token", "", "Token for token based auth protected communication with FHIR Server")
	_ = viper.BindPFlag("retrieve.fhirServerToken", command.PersistentFlags().Lookup("fhir-server-token"))

	command.PersistentFlags().BoolVar(&retrieveOpts.dev, "dev", false, "Mounts main.R, scripts/ and assets/ from current working directory for local development.")

	command.PersistentFlags().StringToStringP("env", "e", map[string]string{}, "Accepts key-value pairs in the form of key=value and passes them unchanged to the running scripts")
	_ = viper.BindPFlag("retrieve.env", command.PersistentFlags().Lookup("env"))

	return command
}
