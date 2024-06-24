package automator

import (
	"context"

	"smuggr.xyz/spammr/common/configurator"
	"smuggr.xyz/spammr/common/logger"

	"github.com/chromedp/chromedp"
)

var Logger = logger.NewCustomLogger("auto")
var Config *configurator.AutomatorConfig

var Automators map[string]Automator

func SetupBrowser(ctx context.Context, config *configurator.AutomatorConfig) (context.Context, context.CancelFunc) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", config.Headless),
	)

	actx, acancel := chromedp.NewExecAllocator(context.Background(), opts...)

	ctx, cancel := chromedp.NewContext(
		actx,
		chromedp.WithDebugf(Logger.Printf),
		chromedp.WithLogf(Logger.Printf),
	)

	return ctx, func() {
		cancel()
		acancel()
	}
}

func Initialize(cmdFlags *configurator.CmdFlags) {
	Logger.Log(logger.MsgInitializing)

	Config = &configurator.Config.Automator

	var err error
	Automators, err = LoadAutomatorFiles()
	if err != nil {
		Logger.Fatal(err)
	}

	for _, automator := range Automators {
		ctx, cancel := SetupBrowser(context.Background(), Config)
		defer cancel()

		if err := RunAutomator(ctx, &automator); err != nil {
			Logger.Fatal(err)
		}
	}
}