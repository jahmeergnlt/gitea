package webhook

import (
	"log"

	"gitea/models/webhook"
)

func Init() error {
	tasks, err := webhook.ListProcessingTasks()
	if err != nil {
		return err
	}
	for _, task := range tasks {
		task.IsProcessing = false
		task.IsDelivered = true
		task.IsSucceed = false
		if task.Retries < 3 {
			task.IsDelivered = false
		}
		err := webhook.UpdateTask(task)
		if err != nil {
			log.Printf("failed to update recovered task %d: %v", task.ID, err)
		}
	}
	return nil
}
