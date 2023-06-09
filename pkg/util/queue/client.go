package queue

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/lactobasilusprotectus/go-template/pkg/common/config"
)

type Interface interface {
	EnqueueTask(task *asynq.Task) (*asynq.TaskInfo, error)
	EnqueueTaskContext(ctx context.Context, task *asynq.Task) (*asynq.TaskInfo, error)
	Close() error
}

// Client is a struct that encapsulates the asynq client.
type Client struct {
	Asynqclient *asynq.Client
}

// NewClient constructs new Client
func NewClient(config config.RedisConfig) *Client {
	opts := asynq.RedisClientOpt{
		Addr:     config.Host,
		Username: config.Username,
		Password: config.Password,
		DB:       config.DB,
	}

	return &Client{Asynqclient: asynq.NewClient(opts)}
}

// EnqueueTask enqueues a task to be processed by the asynq server.
func (c *Client) EnqueueTask(task *asynq.Task) (*asynq.TaskInfo, error) {
	return c.Asynqclient.Enqueue(task)
}

// EnqueueTaskContext enqueues a task to be processed by the asynq server.
func (c *Client) EnqueueTaskContext(ctx context.Context, task *asynq.Task) (*asynq.TaskInfo, error) {
	return c.Asynqclient.EnqueueContext(ctx, task)
}

// Close closes the client connection to redis.
func (c *Client) Close() error {
	return c.Asynqclient.Close()
}
