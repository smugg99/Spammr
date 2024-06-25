package automator

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

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

func runAction(ctx context.Context, action Action) error {
	switch action.Type {
	case ActionNavigate:
		return navigateAction(ctx, action.Selector)
	case ActionWait:
		return waitAction(ctx, action.Selector, action.Duration)
	case ActionFill:
		return fillAction(ctx, action.Selector, fmt.Sprintf("%v", action.Value))
	case ActionReturn:
		return returnAction(action.Value)
	case ActionPrint:
		return printAction(action.Value)
	default:
		return logger.ErrUnknownActionType.Format(action.Type)
	}
}

func executeActions(ctx context.Context, actions []Action) error {
	for index, action := range actions {
		start := time.Now()
		ProgressLogger.Progressf("[%d] %s", index, action.Type)

		if err := runAction(ctx, action); err != nil {
			// duration := time.Since(start)

			// ProgressLogger.ProgressErrorf("[%d] %s failed after %v: %v", index, action.Type, duration, err)
			
			if err := onFailure(ctx, action); err != nil {
				return err
			}

			return err
		}

		duration := time.Since(start)
		ProgressLogger.ProgressDebugf("[%d] %s completed in %v", index, action.Type, duration)
	}
	return nil
}

func RunAutomator(ctx context.Context, automator *Automator) error {
	ProgressLogger.Progressf("[%s]", automator.Name)
	return executeActions(ctx, automator.Actions)
}
