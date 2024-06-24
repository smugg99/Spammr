package commander

import (
	"os"

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

var CmdFlags = &configurator.CmdFlags{}

func Initialize(runFunc func(cmd *cobra.Command, args []string)) {
	RunFunc = runFunc

	SetupVersionCmd()
	SetupRunCmd()

	SetupRootCmd()
}

func ExecuteRootCommand() {
	if err := RootCmd.Execute(); err != nil {
		// Logger.Errorf("error while executing command: %v", err)
		os.Exit(1)
	}
}
