package automator

import (
	"context"
	"fmt"
	"time"

	"smuggr.xyz/spammr/common/logger"
)

func runAction(ctx context.Context, action Action) error {
	switch action.Type {
	case ActionNavigate:
		return navigateAction(ctx, action.Value)
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
			ProgressLogger.ProgressErrorf("[%d] %s : %v [%v]", index, action.Type, err, time.Since(start))

			onFailure(ctx, action)

			return err
		}

		ProgressLogger.ProgressDebugf("[%d] %s [%v]", index, action.Type, time.Since(start))
	}
	return nil
}

func RunAutomator(ctx context.Context, automator *Automator) error {
	ProgressLogger.Progressf("[%s]", automator.Name)
	return executeActions(ctx, automator.Actions)
}
