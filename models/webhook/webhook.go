package webhook

import (
	"errors"
	"sync"
)

type WebhookTask struct {
	ID             int64
	HookID         int64
	UUID           string
	PayloadContent string
	IsDelivered    bool
	IsSucceed      bool
	IsProcessing   bool
	Retries        int
}

var (
	tasks   = make(map[int64]*WebhookTask)
	tasksMu sync.RWMutex
	nextID  int64 = 1
)

func CreateTask(task *WebhookTask) error {
	tasksMu.Lock()
	defer tasksMu.Unlock()
	task.ID = nextID
	nextID++
	tasks[task.ID] = task
	return nil
}

func GetTask(id int64) (*WebhookTask, error) {
	tasksMu.RLock()
	defer tasksMu.RUnlock()
	task, ok := tasks[id]
	if !ok {
		return nil, errors.New("task not found")
	}
	copyTask := *task
	return &copyTask, nil
}

func UpdateTask(task *WebhookTask) error {
	tasksMu.Lock()
	defer tasksMu.Unlock()
	if _, ok := tasks[task.ID]; !ok {
		return errors.New("task not found")
	}
	tasks[task.ID] = task
	return nil
}

func ListProcessingTasks() ([]*WebhookTask, error) {
	tasksMu.RLock()
	defer tasksMu.RUnlock()
	var list []*WebhookTask
	for _, task := range tasks {
		if task.IsProcessing {
			copyTask := *task
			list = append(list, &copyTask)
		}
		}
	return list, nil
}

func ClearTasks() {
	tasksMu.Lock()
	defer tasksMu.Unlock()
	tasks = make(map[int64]*WebhookTask)
	nextID = 1
}
