package automator

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"smuggr.xyz/spammr/common/logger"
)

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

func LoadAutomatorFiles() (map[string]Automator, error) {
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

			Logger.Debugf("automator %s loaded: %v", filename, automator)
		}
	}

	return automators, nil
}

func RunAutomator(ctx context.Context, automator *Automator) error {
	for _, action := range automator.Actions {
		switch action.Type {
		case ActionNavigate:
			if err := navigateAction(ctx, action.Selector); err != nil {
				return err
			}
		case ActionWait:
			if err := waitAction(ctx, action.Selector, action.Duration); err != nil {
				return err
			}
		case ActionFill:
			if err := fillAction(ctx, action.Selector, action.Value); err != nil {
				return err
			}
		default:
			return logger.ErrUnknownActionType.Format(action.Type)
		}
	}
	
	return nil
}