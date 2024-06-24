package main

import (
	"smuggr.xyz/spammr/common/configurator"
	"smuggr.xyz/spammr/common/logger"
	"smuggr.xyz/spammr/core/commander"
	"smuggr.xyz/spammr/core/automator"

	"github.com/spf13/cobra"
)

func main() {
	commander.Initialize(func(cmd *cobra.Command, args []string) {
		commander.Logger.Log(logger.MsgInitializing)
	
		configurator.Initialize()

		logger.Initialize(commander.CmdFlags.IsVerbose)

		automator.Initialize(commander.CmdFlags)
	})
	
	commander.ExecuteRootCommand()
}
