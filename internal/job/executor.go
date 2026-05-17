package job

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"

	"github.com/feishu/feishu-office-suite/internal/data"
)

type Executor struct {
	data   *data.Data
	cron   *cron.Cron
	redis  *redis.Client
}

func NewExecutor(data *data.Data) *Executor {
	return &Executor{
		data:  data,
		cron:  cron.New(),
		redis: data.Redis(),
	}
}

type TaskType string

const (
	TaskTypeSendMessage     TaskType = "send_message"
	TaskTypeSyncUser        TaskType = "sync_user"
	TaskTypeSyncDepartment  TaskType = "sync_department"
	TaskTypeSendNotification TaskType = "send_notification"
	TaskTypeProcessApproval TaskType = "process_approval"
)

type Task struct {
	Type    TaskType     `json:"type"`
	Payload TaskPayload  `json:"payload"`
}

type TaskPayload map[string]interface{}

func (e *Executor) RegisterHandlers() {
	if err := e.cron.AddFunc("0 */5 * * * *", e.syncUsersJob); err != nil {
		log.Printf("Failed to register sync users job: %v", err)
	}
	if err := e.cron.AddFunc("0 */10 * * * *", e.syncDepartmentsJob); err != nil {
		log.Printf("Failed to register sync departments job: %v", err)
	}
	e.cron.Start()
}

func (e *Executor) ExecuteTask(ctx context.Context, t *Task) error {
	switch t.Type {
	case TaskTypeSendMessage:
		return e.sendMessage(ctx, t.Payload)
	case TaskTypeSyncUser:
		return e.syncUser(ctx, t.Payload)
	case TaskTypeSyncDepartment:
		return e.syncDepartment(ctx, t.Payload)
	case TaskTypeSendNotification:
		return e.sendNotification(ctx, t.Payload)
	case TaskTypeProcessApproval:
		return e.processApproval(ctx, t.Payload)
	default:
		return fmt.Errorf("unknown task type: %s", t.Type)
	}
}

func (e *Executor) sendMessage(ctx context.Context, payload TaskPayload) error {
	receiveID, _ := payload["receive_id"].(string)
	msgType, _ := payload["msg_type"].(string)
	content, _ := payload["content"].(string)

	log.Printf("Sending message to %s: type=%s, content=%s", receiveID, msgType, content)
	return nil
}

func (e *Executor) syncUsersJob() {
	log.Println("Running scheduled user sync job")
	task := &Task{
		Type:    TaskTypeSyncUser,
		Payload: TaskPayload{},
	}
	if err := e.ExecuteTask(context.Background(), task); err != nil {
		log.Printf("Failed to sync users: %v", err)
	}
}

func (e *Executor) syncDepartmentsJob() {
	log.Println("Running scheduled department sync job")
	task := &Task{
		Type:    TaskTypeSyncDepartment,
		Payload: TaskPayload{},
	}
	if err := e.ExecuteTask(context.Background(), task); err != nil {
		log.Printf("Failed to sync departments: %v", err)
	}
}

func (e *Executor) syncUser(ctx context.Context, payload TaskPayload) error {
	log.Println("Syncing user data from Feishu")
	return nil
}

func (e *Executor) syncDepartment(ctx context.Context, payload TaskPayload) error {
	log.Println("Syncing department data from Feishu")
	return nil
}

func (e *Executor) sendNotification(ctx context.Context, payload TaskPayload) error {
	log.Println("Sending notification")
	return nil
}

func (e *Executor) processApproval(ctx context.Context, payload TaskPayload) error {
	log.Println("Processing approval")
	return nil
}

type AsynqWorker struct {
	executor *Executor
	redisURL string
}

func NewAsynqWorker(executor *Executor) *AsynqWorker {
	return &AsynqWorker{
		executor: executor,
	}
}

func (w *AsynqWorker) Start() error {
	w.executor.RegisterHandlers()

	cfg := sarama.NewConfig()
	cfg.ClientID = "feishu-asynq-worker"

	log.Println("Asynq worker started")
	select {}
}

func (w *AsynqWorker) Enqueue(task *Task) error {
	taskJSON, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	return w.executor.redis.LPush(context.Background(), "asynq:pending", taskJSON).Err()
}