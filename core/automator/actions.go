package automator

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
)

func navigateAction(ctx context.Context, url string) error {
	return chromedp.Run(ctx, chromedp.Navigate(url))
}

func waitAction(ctx context.Context, selector string, duration int) error {
	if err := chromedp.Run(ctx, chromedp.WaitVisible(selector)); err != nil {
		return err
	}
	time.Sleep(time.Duration(duration) * time.Millisecond)
	return nil
}

func fillAction(ctx context.Context, selector string, value string) error {
	return chromedp.Run(ctx, chromedp.SetValue(selector, value))
}