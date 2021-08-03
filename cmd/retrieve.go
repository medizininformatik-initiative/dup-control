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

type RetrieveOpts struct {
	Workpackage        string //req
	Site               string //req
	fhirServerEndpoint string
	fhirServerUser     string
	fhirServerPass     string
	fhirServerCACert   string
}

var retrieveOpts = RetrieveOpts{}

func createRetrieveOpts(retrieveOpts RetrieveOpts) (container.PullOpts, container.RunOpts) {
	pullOpts := container.PullOpts{
		Workpackage: retrieveOpts.Workpackage,
		Site:        retrieveOpts.Site,
	}
	runOpts := container.RunOpts{
		Env: []string{
			fmt.Sprintf("FHIR_SERVER_ENDPOINT=\"%s\"", retrieveOpts.fhirServerEndpoint),
		},
		Mounts: []docker.Mount{
			container.LocalMount("outputLocal", true),
			container.LocalMount("outputGlobal", true),
		},
	}
	if runtime.GOOS != "windows" {
		runOpts.User = fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid())
	}
	if retrieveOpts.fhirServerUser != "" && retrieveOpts.fhirServerPass != "" {
		runOpts.Env = append(runOpts.Env,
			fmt.Sprintf("FHIR_SERVER_USER=\"%s\"", retrieveOpts.fhirServerUser),
			fmt.Sprintf("FHIR_SERVER_PASS=\"%s\"", retrieveOpts.fhirServerPass))
	}
	if retrieveOpts.fhirServerCACert != "" {
		runOpts.Mounts = append(runOpts.Mounts,
			docker.Mount{
				Source:      retrieveOpts.fhirServerCACert,
				Destination: "/etc/ssl/certs/ca-certificates.crt",
				Driver:      "local",
				RW:          false})
	}
	return pullOpts, runOpts
}

var retrieveCommand = &cobra.Command{
	Use:   "retrieve",
	Short: "Retrieve bundles from FHIR server",
	Long:  "You can retrieve bundles from the FHIR server for a specific POLAR workpackage",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if viper.GetString("fhirServerEndpoint") == "" {
			return fmt.Errorf("fhirServerEndpoint not set")
		} else {
			retrieveOpts.fhirServerEndpoint = viper.GetString("fhirServerEndpoint")
		}
		retrieveOpts.Site = viper.GetString("site")
		retrieveOpts.fhirServerUser = viper.GetString("fhirServerUser")
		retrieveOpts.fhirServerPass = viper.GetString("fhirServerPass")
		retrieveOpts.fhirServerCACert = viper.GetString("fhirServerCACert")
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

	retrieveCommand.PersistentFlags().StringVar(&retrieveOpts.Workpackage, "wp", "", "Workpackage to execute (e.g. 'wp-1-1-pilot').")
	_ = retrieveCommand.MarkPersistentFlagRequired("wp")

	retrieveCommand.PersistentFlags().String("site", "latest", "Determines which image to use, as images are (not necessarily) hand-tailored for different dic sites. (e.g. 'dic-giessen', 'dic-leipzig', 'dic-muenchen').")
	_ = viper.BindPFlag("site", retrieveCommand.PersistentFlags().Lookup("site"))

	retrieveCommand.PersistentFlags().String("fhirServerEndpoint", "", "the base URL of the FHIR server to use")
	_ = viper.BindPFlag("fhirServerEndpoint", retrieveCommand.PersistentFlags().Lookup("fhirServerEndpoint"))

	retrieveCommand.PersistentFlags().String("fhirServerUser", "", "fhirServerUser for basic auth protected communication with FHIR server")
	_ = viper.BindPFlag("fhirServerUser", retrieveCommand.PersistentFlags().Lookup("fhirServerUser"))

	retrieveCommand.PersistentFlags().String("fhirServerPass", "", "fhirServerPass for basic auth protected communication with FHIR server")
	_ = viper.BindPFlag("fhirServerPass", retrieveCommand.PersistentFlags().Lookup("fhirServerPass"))

	retrieveCommand.PersistentFlags().String("fhirServerCACert", "", "CA Certificate file for https connection to FHIR Server")
	_ = viper.BindPFlag("fhirServerCACert", retrieveCommand.PersistentFlags().Lookup("fhirServerCACert"))
}
