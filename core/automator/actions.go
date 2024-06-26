package automator

import (
	"bufio"
	"context"
	"os"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"smuggr.xyz/spammr/common/logger"
)

func navigateAction(ctx context.Context, url interface{}) error {
	switch v := url.(type) {
	case string:
		return chromedp.Run(ctx, chromedp.Navigate(v))
	}

	return logger.ErrUnsupportedActionValueType
}

func waitAction(ctx context.Context, selector string, duration int) error {
	if selector != "" {
		var waitCtx context.Context
		var cancel context.CancelFunc

		if duration > 0 {
			waitCtx, cancel = context.WithTimeout(ctx, time.Duration(duration)*time.Millisecond)
		} else {
			waitCtx, cancel = context.WithCancel(ctx)
		}
		defer cancel()

		return chromedp.Run(waitCtx, chromedp.WaitVisible(selector))
	} else if duration > 0 {
		time.Sleep(time.Duration(duration) * time.Millisecond)
		return nil
	}
	return nil
}

func fillAction(ctx context.Context, selector string, value string) error {
	return chromedp.Run(ctx, chromedp.SetValue(selector, value))
}

func returnAction(value interface{}) error {
	switch v := value.(type) {
	case bool:
		if !v {
			return logger.ErrActionReturnedFalse
		}
	case string:
		if v != "true" {
			return logger.ErrActionReturnedFalse
		}
	}

	return logger.ErrUnsupportedActionValueType
}

func printAction(value interface{}) error {
	ProgressLogger.Progress(value)
	return nil
}

func promptConfirm() bool {
	ProgressLogger.Progress(logger.MsgEnterYesOrNo)

	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			ProgressLogger.ProgressError("Error reading input:", err)
			return false
		}

		input = strings.TrimSpace(strings.ToUpper(input))

		switch input {
		case "Y", "YES":
			return true
		case "N", "NO":
			return false
		default:
			ProgressLogger.ProgressError(logger.ErrInvalidYesOrNoInput)
		}
	}
}

func confirmAction() error {
	if promptConfirm() {
		return nil
	}

	return logger.ErrUserChoseToExit
}

func onFailure(ctx context.Context, action Action) error {
	if len(action.OnFailure) > 0 {
		return executeActions(ctx, action.OnFailure)
	}

	return nil
}