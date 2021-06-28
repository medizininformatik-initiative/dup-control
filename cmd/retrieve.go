package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var fhirServerEndpoint string
var username string
var password string
var cacert string

var retrieveCommand = &cobra.Command{
	Use:   "retrieve",
	Short: "Retrieve bundles from FHIR server",
	Long:  "You can retrieve bundles from the FHIR server for a specific POLAR workpackage",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("FHIR server is %s", fhirServerEndpoint)
	},
}

func init() {
	rootCmd.AddCommand(retrieveCommand)

	retrieveCommand.PersistentFlags().StringVar(&fhirServerEndpoint, "fhirServerEndpoint", "", "the base URL of the FHIR server to use")
	_ = vip.BindPFlag("fhirServerEndpoint", retrieveCommand.PersistentFlags().Lookup("fhirServerEndpoint"))
	//_ = retrieveCommand.MarkPersistentFlagRequired("fhirServerEndpoint")

	retrieveCommand.PersistentFlags().StringVar(&username, "username", "", "username for basic auth protected communication with FHIR server")
	_ = vip.BindPFlag("username", retrieveCommand.PersistentFlags().Lookup("username"))

	retrieveCommand.PersistentFlags().StringVar(&password, "password", "", "password for basic auth protected communication with FHIR server")
	_ = vip.BindPFlag("password", retrieveCommand.PersistentFlags().Lookup("password"))

	retrieveCommand.PersistentFlags().StringVar(&cacert, "cacert", "", "CA Certificate file for https connection to FHIR Server")
	_ = vip.BindPFlag("cacert", retrieveCommand.PersistentFlags().Lookup("cacert"))
}
