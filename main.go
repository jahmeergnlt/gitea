package main

import (
	"context"
	"fmt"
	"time"

	"gitea/models/webhook"
	services "gitea/services/webhook"
)

func main() {
	fmt.Println("Starting Webhook Service Demo...")

	// 1. Simulate Orphaned Task Recovery
	fmt.Println("\n--- Scenario 1: Orphaned Task Recovery ---")
	task1 := &webhook.WebhookTask{
		UUID:         "orphaned-task-1",
		IsProcessing: true,
		Retries:      1,
	}
	_ = webhook.CreateTask(task1)
	fmt.Printf("Created orphaned task: ID=%d, IsProcessing=%t, IsDelivered=%t, Retries=%d\n",
		task1.ID, task1.IsProcessing, task1.IsDelivered, task1.Retries)

	fmt.Println("Initializing Webhook Service (recovering tasks)...")
	_ = services.Init()

	recoveredTask1, _ := webhook.GetTask(task1.ID)
	fmt.Printf("After recovery: ID=%d, IsProcessing=%t, IsDelivered=%t, Retries=%d\n",
		recoveredTask1.ID, recoveredTask1.IsProcessing, recoveredTask1.IsDelivered, recoveredTask1.Retries)

	// 2. Simulate Graceful Shutdown Handling
	fmt.Println("\n--- Scenario 2: Graceful Shutdown Handling ---")
	task2 := &webhook.WebhookTask{
		UUID: "shutdown-task-1",
	}
	_ = webhook.CreateTask(task2)
	fmt.Printf("Created task: ID=%d, IsProcessing=%t, IsDelivered=%t\n",
		task2.ID, task2.IsProcessing, task2.IsDelivered)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(50 * time.Millisecond)
		fmt.Println("Simulating graceful shutdown (canceling context)...")
		cancel()
	}()

	fmt.Println("Delivering task...")
	err := services.Deliver(ctx, task2)
	fmt.Printf("Delivery finished with error: %v\n", err)

	finalTask2, _ := webhook.GetTask(task2.ID)
	fmt.Printf("Final task state: ID=%d, IsProcessing=%t, IsDelivered=%t, IsSucceed=%t\n",
		finalTask2.ID, finalTask2.IsProcessing, finalTask2.IsDelivered, finalTask2.IsSucceed)
}
