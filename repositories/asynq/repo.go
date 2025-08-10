package asynq

import (
	"NotificationManagement/config"
	"NotificationManagement/logger"
	"NotificationManagement/types"
	"encoding/json"
	"errors"
	"time"

	"github.com/hibiken/asynq"
)

type Repository struct {
	client    *asynq.Client
	inspector *asynq.Inspector
}

func NewRepository(client *asynq.Client, inspector *asynq.Inspector) *Repository {
	return &Repository{
		client:    client,
		inspector: inspector,
	}
}

func (repo *Repository) CreateTask(event types.AsynqTaskType, data interface{}) (*asynq.Task, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(event.String(), payload), nil
}

func (repo *Repository) EnqueueTask(task *asynq.Task, customOpts *types.AsynqOption) (string, error) {
	opts := repo.asynqOptions(*customOpts)
	taskInfo, err := repo.client.Enqueue(task, opts...)
	if err != nil {
		return "", err
	}
	return taskInfo.ID, nil
}

func (repo *Repository) DequeueTask(taskID string) error {
	existingTask, err := repo.inspector.GetTaskInfo(config.Asynq().Queue, taskID)
	if err != nil && !errors.Is(err, asynq.ErrTaskNotFound) {
		return err
	}
	if existingTask == nil {
		return nil
	}

	deleteOrCancelTask := func(task *asynq.TaskInfo) error {
		if task.State != asynq.TaskStateActive {
			repo.inspector.DeleteTask(config.Asynq().Queue, task.ID)
		}
		if err := repo.inspector.CancelProcessing(task.ID); err != nil {
			return err
		}
		return repo.inspector.DeleteTask(config.Asynq().Queue, task.ID)
	}

	err = deleteOrCancelTask(existingTask)
	if errors.Is(err, asynq.ErrTaskNotFound) || errors.Is(err, asynq.ErrQueueNotFound) {
		return nil
	}
	if err != nil {
		logger.Error("error on deleting task ", taskID, " : ", err)
		return err
	}

	return nil
}

func (repo *Repository) asynqOptions(customOpts types.AsynqOption) []asynq.Option {
	retryOpt := asynq.MaxRetry(0)
	queueOpt := asynq.Queue(config.Asynq().Queue)
	retentionOpt := asynq.Retention(time.Duration(*config.Asynq().Retention) * time.Hour)

	if customOpts.Retry > 0 {
		retryOpt = asynq.MaxRetry(customOpts.Retry)
	}

	if customOpts.Queue != "" {
		queueOpt = asynq.Queue(customOpts.Queue)
	}

	if customOpts.RetentionHours > 0 {
		retentionOpt = asynq.Retention(customOpts.RetentionHours * time.Hour)
	}

	opts := []asynq.Option{
		retryOpt,
		queueOpt,
		retentionOpt,
	}

	// zero value not allowed
	if len(customOpts.TaskID) > 0 {
		opts = append(opts, asynq.TaskID(customOpts.TaskID))
	}

	// zero value not allowed
	if customOpts.DelaySeconds > 0 {
		opts = append(opts, asynq.ProcessIn(customOpts.DelaySeconds*time.Second))
	}

	// zero value not allowed
	if customOpts.UniqueTTLSeconds > 0 {
		opts = append(opts, asynq.Unique(customOpts.UniqueTTLSeconds*time.Second))
	}

	return opts
}
