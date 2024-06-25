package automator

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"smuggr.xyz/spammr/common/configurator"
	"smuggr.xyz/spammr/common/logger"

	"github.com/chromedp/chromedp"
)

var Logger = logger.NewCustomLogger("auto")
var ProgressLogger = logger.NewCustomLogger("")

var Config *configurator.AutomatorConfig

var Automators map[string]Automator

func readAutomatorFromFile(filePath string) (Automator, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return Automator{}, err
	}

	var automator Automator
	if err := json.Unmarshal(data, &automator); err != nil {
		return Automator{}, err
	}

	return automator, nil
}

func SetupBrowser() (context.Context, context.CancelFunc) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", Config.Headless),
	)

	actx, acancel := chromedp.NewExecAllocator(context.Background(), opts...)

	ctx, cancel := chromedp.NewContext(
		actx,
		chromedp.WithDebugf(func(format string, args ...interface{}) {
			if Config.AttachDebug {
				Logger.Printf(format, args...)
			}
		}),
		chromedp.WithLogf(func(format string, args ...interface{}) {
			if Config.AttachLog {
				Logger.Printf(format, args...)
			}
		}),
	)

	return ctx, func() {
		cancel()
		acancel()
	}
}

func LoadAutomatorFiles(cmdFlags *configurator.CmdFlags) (map[string]Automator, error) {
	_directory := "AUTOMATORS_DIRECTORY"
	directory := os.Getenv(_directory)
	if directory == "" {
		return nil, logger.ErrEnvVariableNotSet.Format(_directory)
	}

	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	_automators_ext := "AUTOMATOR_FILE_EXTENSION"
	automators_ext := os.Getenv(_automators_ext)
	if directory == "" {
		return nil, logger.ErrEnvVariableNotSet.Format(_directory)
	}

	automators := make(map[string]Automator)

	for _, file := range files {
		filename := file.Name()
		if !file.IsDir() && filepath.Ext(filename) == automators_ext {
			filePath := filepath.Join(directory, filename)

			Logger.Info(logger.MsgLoadingResource.Format(filePath, logger.ResourceAutomator))

			automator, err := readAutomatorFromFile(filePath)
			if err != nil {
				Logger.Error(logger.ErrReadingResource.Format(filePath, logger.ResourceAutomator))
				return nil, err
			}

			automators[filename] = automator

			Logger.Log(logger.MsgResourceLoaded.Format(filePath, logger.ResourceAutomator))
		}
	}

	return automators, nil
}

func Initialize(cmdFlags *configurator.CmdFlags) {
	Logger.Log(logger.MsgInitializing)

	Config = &configurator.Config.Automator

	var err error
	Automators, err = LoadAutomatorFiles(cmdFlags)
	if err != nil {
		Logger.Fatal(err)
	}

	for _, automator := range Automators {
		ctx, cancel := SetupBrowser()
		defer cancel()

		ReplacePlaceholders(&automator, cmdFlags)

		if err := RunAutomator(ctx, &automator); err != nil {
			Logger.Log(logger.ErrAutomatorError.Format(automator.Name, err))
		} else {
			Logger.Log(logger.MsgAutomatorSuccess.Format(automator.Name))
		}
	}
}