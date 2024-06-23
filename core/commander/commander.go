package commander

import (
	"fmt"
	"os"

	"smuggr.xyz/spammr/common/logger"

	"github.com/spf13/cobra"
)

var Logger *logger.CustomLogger
var RootCmd *cobra.Command

func Initialize(runFunc func(cmd *cobra.Command, args []string)) {
	var (
        verbose bool
    )

	RootCmd = &cobra.Command{
		Use:   "spammr",
		Short: "Spammr is a simple CLI application built to aid with spamming tasks.",
		Long:  `Spammr is a simple CLI application built to automate spam data submission on various sites that accept forms like phone numbers, emails, addresses, and more.`,
		Run: runFunc,
	}

    RootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose mode")

	RootCmd.AddCommand(&cobra.Command{
        Use:   "version",
        Short: "Prints the version of Spammr",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("Spammr v0.1.0")
        },
    })

	RootCmd.AddCommand(&cobra.Command{
        Use:   "run",
        Short: "Runs the Spammr application",
        Run: func(cmd *cobra.Command, args []string) {
            logger.Initialize(verbose)
			
			Logger = logger.NewCustomLogger("cmdr")
			Logger.Log(logger.MsgInitializing)
        },
    })
}

func ExecuteRootCommand() {
	if err := RootCmd.Execute(); err != nil {
		Logger.Errorf("error while executing command: %v", err)
		os.Exit(1)
	}
}
