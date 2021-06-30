package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fhirServerUser string
var fhirServerPass string
var cacert string
var fhirServerEndpoint string

var retrieveCommand = &cobra.Command{
	Use:   "retrieve",
	Short: "Retrieve bundles from FHIR server",
	Long:  "You can retrieve bundles from the FHIR server for a specific POLAR workpackage",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if viper.GetString("fhirServerEndpoint") == "" {
			return fmt.Errorf("fhirServerEndpoint not set")
		} else {
			fhirServerEndpoint = viper.GetString("fhirServerEndpoint")
		}
		fhirServerUser = viper.GetString("fhirServerUser")
		fhirServerPass = viper.GetString("fhirServerPass")
		cacert = viper.GetString("cacert")
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := containerRuntime.Pull(workpackage, site); err != nil {
			return err
		}

		if err := containerRuntime.Run("retrieval", workpackage, site); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(retrieveCommand)

	retrieveCommand.PersistentFlags().String("fhirServerEndpoint", "", "the base URL of the FHIR server to use")
	_ = viper.BindPFlag("fhirServerEndpoint", retrieveCommand.PersistentFlags().Lookup("fhirServerEndpoint"))

	retrieveCommand.PersistentFlags().String("fhirServerUser", "", "fhirServerUser for basic auth protected communication with FHIR server")
	_ = viper.BindPFlag("fhirServerUser", retrieveCommand.PersistentFlags().Lookup("fhirServerUser"))

	retrieveCommand.PersistentFlags().String("fhirServerPass", "", "fhirServerPass for basic auth protected communication with FHIR server")
	_ = viper.BindPFlag("fhirServerPass", retrieveCommand.PersistentFlags().Lookup("fhirServerPass"))

	retrieveCommand.PersistentFlags().String("cacert", "", "CA Certificate file for https connection to FHIR Server")
	_ = viper.BindPFlag("cacert", retrieveCommand.PersistentFlags().Lookup("cacert"))
}
