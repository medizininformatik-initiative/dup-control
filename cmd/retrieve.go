package cmd

import (
	"fmt"
	coll "git.smith.care/smith/uc-phep/polar/polarctl/util"
	"git.smith.care/smith/uc-phep/polar/polarctl/util/container"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"runtime"
)

type RetrieveOpts struct {
	workpackage        string //req
	site               string //req
	fhirServerEndpoint string
	fhirServerUser     string
	fhirServerPass     string
	fhirServerCACert   string
	dev                bool
	env                map[string]string
}

var retrieveOpts = RetrieveOpts{}

func createRetrieveOpts(retrieveOpts RetrieveOpts) (container.PullOpts, container.RunOpts) {
	pullOpts := container.PullOpts{
		Image: retrieveOpts.workpackage,
		Tag:   retrieveOpts.site,
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

var retrieveCommand = &cobra.Command{
	Use:   "retrieve",
	Short: "Retrieve bundles from FHIR server",
	Long:  "You can retrieve bundles from the FHIR server for a specific POLAR workpackage",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if viper.GetString("retrieve.fhirServerEndpoint") == "" {
			return fmt.Errorf("retrieve.fhirServerEndpoint not set")
		} else {
			retrieveOpts.fhirServerEndpoint = viper.GetString("retrieve.fhirServerEndpoint")
		}
		retrieveOpts.site = viper.GetString("retrieve.site")
		retrieveOpts.fhirServerUser = viper.GetString("retrieve.fhirServerUser")
		retrieveOpts.fhirServerPass = viper.GetString("retrieve.fhirServerPass")
		retrieveOpts.fhirServerCACert = viper.GetString("retrieve.fhirServerCACert")
		retrieveOpts.env = viper.GetStringMapString("retrieve.env")
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		pullOpts, runOpts := createRetrieveOpts(retrieveOpts)
		if err := containerRuntime.Pull(pullOpts); err != nil {
			return err
		}

		if err := containerRuntime.Run("retrieval", pullOpts, runOpts); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(retrieveCommand)

	retrieveCommand.PersistentFlags().StringVar(&retrieveOpts.workpackage, "wp", "", "Image to execute (e.g. 'wp-1-1-pilot').")
	_ = retrieveCommand.MarkPersistentFlagRequired("wp")

	retrieveCommand.PersistentFlags().String("site", "latest", "Determines which image to use, as images are (not necessarily) hand-tailored for different dic sites. (e.g. 'dic-giessen', 'dic-leipzig', 'dic-muenchen').")
	_ = viper.BindPFlag("retrieve.site", retrieveCommand.PersistentFlags().Lookup("site"))

	retrieveCommand.PersistentFlags().String("fhir-server-endpoint", "", "the base URL of the FHIR server to use")
	_ = viper.BindPFlag("retrieve.fhirServerEndpoint", retrieveCommand.PersistentFlags().Lookup("fhir-server-endpoint"))

	retrieveCommand.PersistentFlags().String("fhir-server-user", "", "fhirServerUser for basic auth protected communication with FHIR server")
	_ = viper.BindPFlag("retrieve.fhirServerUser", retrieveCommand.PersistentFlags().Lookup("fhir-server-user"))

	retrieveCommand.PersistentFlags().String("fhir-server-pass", "", "fhirServerPass for basic auth protected communication with FHIR server")
	_ = viper.BindPFlag("retrieve.fhirServerPass", retrieveCommand.PersistentFlags().Lookup("fhir-server-pass"))

	retrieveCommand.PersistentFlags().String("fhir-server-cacert", "", "CA Certificate file for https connection to FHIR Server")
	_ = viper.BindPFlag("retrieve.fhirServerCACert", retrieveCommand.PersistentFlags().Lookup("fhir-server-cacert"))

	retrieveCommand.PersistentFlags().BoolVar(&retrieveOpts.dev, "dev", false, "Mounts main.R, scripts/ and assets/ from current working directory for local development.")

	retrieveCommand.PersistentFlags().StringToStringP("env", "e", map[string]string{}, "Accepts key-value pairs in the form of key=value and passes them unchanged to the running scripts")
	_ = viper.BindPFlag("retrieve.env", retrieveCommand.PersistentFlags().Lookup("env"))
}
