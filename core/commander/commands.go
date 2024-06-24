package commander

import (
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

func SetupRootCmd() {
	RootCmd = &cobra.Command{
		Use:   "spammr",
		Short: "Spammr is a simple CLI application built to aid with spamming tasks.",
		Long:  `Spammr is a simple CLI application built to automate spam data submission on various sites that accept forms like phone numbers, emails, addresses, and more.`,
	}

	RootCmd.AddCommand(VersionCmd, RunCmd)
}

func SetupVersionCmd() {
	VersionCmd = &cobra.Command{
		Use:   "version",
		Short: "Prints the version of Spammr",
		Run: func(cmd *cobra.Command, args []string) {
			Logger.Info("Spammr v0.1.0")
		},
	}
}

func SetupRunCmd() {
	RunCmd = &cobra.Command{
		Use:   "run",
		Short: "Runs the Spammr application",
		Run:   RunFunc,
	}

	verboseDefault, err := strconv.ParseBool(os.Getenv("VERBOSE"))
	if err != nil {
		verboseDefault = false
	}

	RunCmd.PersistentFlags().BoolVarP(&CmdFlags.IsVerbose, "verbose", "v", verboseDefault, "Enable verbose mode")
	RunCmd.Flags().StringVarP(&CmdFlags.RequestsDir, "requests-dir", "d", os.Getenv("REQUEST_TEMPLATES_DIRECTORY"), "Directory path for request files")

	RunCmd.Flags().StringVarP(&CmdFlags.Want.Hash, "hash", "x", "", "Specify the desired hash (10 character long random string e.g. 'a450a03cb5')")
	RunCmd.Flags().StringVarP(&CmdFlags.Want.Boundary, "boundary", "b", "", "Specify the desired boundary (16 character long random string e.g. '7MA4YsWxkTrZu0gW')")
	RunCmd.Flags().StringVarP(&CmdFlags.Want.SessionID, "session-id", "s", "", "Specify the desired session ID (10 character long random string e.g. '1sB4XL6kS9')")
	RunCmd.Flags().StringVarP(&CmdFlags.Want.Name, "name", "n", "", "Specify the desired name (first and last, e.g. 'John Doe')")
	RunCmd.Flags().StringVarP(&CmdFlags.Want.FirstName, "first-name", "f", "", "Specify the desired first name (e.g. 'John')")
	RunCmd.Flags().StringVarP(&CmdFlags.Want.LastName, "last-name", "l", "", "Specify the desired last name (e.g. 'Doe')")
	RunCmd.Flags().StringVarP(&CmdFlags.Want.PhoneNumber, "phone-number", "p", "", "Specify the desired phone number (without country code or special characters, e.g. '1234567890')")
	RunCmd.Flags().StringVarP(&CmdFlags.Want.Address, "address", "a", "", "Specify the desired address (street, city, state, zip code)")
	RunCmd.Flags().StringVarP(&CmdFlags.Want.Email, "email", "e", "", "Specify the desired email address")

	RunCmd.Flags().BoolP("help", "h", false, "Help message for run command")
    RunCmd.PersistentFlags().BoolP("help", "", false, "")
}