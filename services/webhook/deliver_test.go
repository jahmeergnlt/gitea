package webhook

import (
	"context"
	"testing"
	"time"

	"gitea/models/webhook"
)

func TestWebhookDelivery_GracefulShutdown(t *testing.T) {
	webhook.ClearTasks()

	task := &webhook.WebhookTask{
		UUID:           "test-uuid-1",
		PayloadContent: "test-payload",
	}
	err := webhook.CreateTask(task)
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	err = Deliver(ctx, task)
	if err == nil {
		t.Fatal("expected context canceled error, got nil")
	}

	updatedTask, err := webhook.GetTask(task.ID)
	if err != nil {
		t.Fatalf("failed to get task: %v", err)
	}

	if updatedTask.IsProcessing {
		t.Error("task should not be in processing state after cancellation")
	}
	if !updatedTask.IsDelivered {
		t.Error("task should be marked as delivered (failed/interrupted) to avoid hanging")
	}
	if updatedTask.IsSucceed {
		t.Error("task should not be marked as succeed")
	}
}

func TestWebhookDelivery_OrphanedTaskRecovery(t *testing.T) {
	webhook.ClearTasks()

	task := &webhook.WebhookTask{
		UUID:         "test-uuid-2",
		IsProcessing: true,
		Retries:      1,
	}
	err := webhook.CreateTask(task)
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}

	err = Init()
	if err != nil {
		t.Fatalf("failed to run Init: %v", err)
	}

	recoveredTask, err := webhook.GetTask(task.ID)
	if err != nil {
		t.Fatalf("failed to get task: %v", err)
	}

	if recoveredTask.IsProcessing {
		t.Error("recovered task should not be in processing state")
	}
	if recoveredTask.IsDelivered {
		t.Error("recovered task with retries < maxRetries should be rescheduled (IsDelivered = false)")
	}
}

func TestWebhookDelivery_OrphanedTaskRecovery_MaxRetries(t *testing.T) {
	webhook.ClearTasks()

	task := &webhook.WebhookTask{
		UUID:         "test-uuid-3",
		IsProcessing: true,
		Retries:      3,
	}
	err := webhook.CreateTask(task)
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}

	err = Init()
	if err != nil {
		t.Fatalf("failed to run Init: %v", err)
	}

	recoveredTask, err := webhook.GetTask(task.ID)
	if err != nil {
		t.Fatalf("failed to get task: %v", err)
	}

	if recoveredTask.IsProcessing {
		t.Error("recovered task should not be in processing state")
	}
	if !recoveredTask.IsDelivered {
		t.Error("recovered task with retries >= maxRetries should be marked as delivered (failed)")
	}
	if recoveredTask.IsSucceed {
		t.Error("recovered task should not be marked as succeed")
	}
}
