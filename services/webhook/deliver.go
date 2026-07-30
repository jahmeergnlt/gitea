package webhook

import (
	"context"
	"errors"
	"time"

	"gitea/models/webhook"
)

func Deliver(ctx context.Context, task *webhook.WebhookTask) error {
	task.IsProcessing = true
	if err := webhook.UpdateTask(task); err != nil {
		return err
	}

	err := doDelivery(ctx, task)

	task.IsProcessing = false
	if err != nil {
		if errors.Is(err, context.Canceled) {
			task.IsDelivered = true
			task.IsSucceed = false
			_ = webhook.UpdateTask(task)
			return err
		}

		task.Retries++
		if task.Retries >= 3 {
			task.IsDelivered = true
			task.IsSucceed = false
		} else {
			task.IsDelivered = false
		}
		_ = webhook.UpdateTask(task)
		return err
	}

	task.IsDelivered = true
	task.IsSucceed = true
	_ = webhook.UpdateTask(task)
	return nil
}

func doDelivery(ctx context.Context, task *webhook.WebhookTask) error {
	select {
	case <-time.After(100 * time.Millisecond):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
