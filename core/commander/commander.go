package commander

import (
	"fmt"
	"os"
	"strconv"

	"smuggr.xyz/spammr/common/configurator"
	"smuggr.xyz/spammr/common/logger"

	"github.com/spf13/cobra"
)

var Logger = logger.NewCustomLogger("cmdr")
var RunFunc func(cmd *cobra.Command, args []string)

var (
	RootCmd    *cobra.Command
	VersionCmd *cobra.Command
	RunCmd     *cobra.Command
)

var CmdFlags configurator.CmdFlags

func Initialize(runFunc func(cmd *cobra.Command, args []string)) {
	RootCmd = &cobra.Command{
		Use:   "spammr",
		Short: "Spammr is a simple CLI application built to aid with spamming tasks.",
		Long:  `Spammr is a simple CLI application built to automate spam data submission on various sites that accept forms like phone numbers, emails, addresses, and more.`,
	}

	VersionCmd = &cobra.Command{
		Use:   "version",
		Short: "Prints the version of Spammr",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Spammr v0.1.0")
		},
	}

	RunCmd = &cobra.Command{
		Use:   "run",
		Short: "Runs the Spammr application",
		Run:   runFunc,
	}

	RunFunc = runFunc
	RootCmd.AddCommand(VersionCmd, RunCmd)

	verboseDefault, err := strconv.ParseBool(os.Getenv("VERBOSE"))
	if err != nil {
		verboseDefault = false
	}

	RunCmd.Flags().BoolVarP(&CmdFlags.IsVerbose, "verbose", "v", verboseDefault, "Enable verbose mode")
	RunCmd.Flags().StringVarP(&CmdFlags.RequestsDir, "requests-dir", "d", os.Getenv("REQUEST_TEMPLATES_DIRECTORY"), "Directory path for request files")

	RunCmd.Flags().StringVarP(&CmdFlags.Want.Boundary, "boundary", "b", "", "Specify the desired boundary")
	RunCmd.Flags().StringVarP(&CmdFlags.Want.SessionID, "session-id", "s", "", "Specify the desired session ID")
	RunCmd.Flags().StringVarP(&CmdFlags.Want.Name, "name", "n", "", "Specify the desired name (first and last)")
	RunCmd.Flags().StringVarP(&CmdFlags.Want.FirstName, "first-name", "f", "", "Specify the desired first name")
	RunCmd.Flags().StringVarP(&CmdFlags.Want.LastName, "last-name", "l", "", "Specify the desired last name")
	RunCmd.Flags().StringVarP(&CmdFlags.Want.PhoneNumber, "phone-number", "p", "", "Specify the desired phone number (without country code)")
	RunCmd.Flags().StringVarP(&CmdFlags.Want.Address, "address", "a", "", "Specify the desired address (street, city, state, zip code)")
	RunCmd.Flags().StringVarP(&CmdFlags.Want.Email, "email", "e", "", "Specify the desired email")
}

func ExecuteRootCommand() {
	if err := RootCmd.Execute(); err != nil {
		// Logger.Errorf("error while executing command: %v", err)
		os.Exit(1)
	}
}
