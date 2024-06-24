package automator

import (
	"smuggr.xyz/spammr/common/logger"
	"smuggr.xyz/spammr/common/configurator"
)

var Logger = logger.NewCustomLogger("auto")

func Initialize(cmdFlags *configurator.CmdFlags) {
	Logger.Log(logger.MsgInitializing)


}