package queue

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/lactobasilusprotectus/go-template/pkg/common/config"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type QueueFunc func(context.Context, *asynq.Task) error

// AsynqServer is a struct that encapsulates the asynq server.
type AsynqServer struct {
	Server    *asynq.Server
	ServerMux *asynq.ServeMux
}

// AsynqServerHandler  wraps queue handler
type AsynqServerHandler interface {
	RegisterQueue(*AsynqServer)
}

// NewAsynqServer constructs new AsynqServer
func NewAsynqServer(config config.RedisConfig) *AsynqServer {
	opts := asynq.RedisClientOpt{
		Addr:     config.Host,
		Username: config.Username,
		Password: config.Password,
		DB:       config.DB,
	}

	server := asynq.NewServer(
		opts,
		asynq.Config{
			Concurrency: 100,
		},
	)

	serverMux := asynq.NewServeMux()

	return &AsynqServer{
		Server:    server,
		ServerMux: serverMux,
	}
}

// AddHandlerFunc adds a custom function to the ServerMux using mux.HandleFunc.
func (as *AsynqServer) AddHandlerFunc(pattern string, handlerFunc QueueFunc) {
	wrappedHandlerFunc := as.WrapHandlerFunc(pattern, handlerFunc)
	as.ServerMux.HandleFunc(pattern, wrappedHandlerFunc)
	log.Println("Added handler func for pattern", pattern)
}

// WrapHandlerFunc wraps queue func
func (as *AsynqServer) WrapHandlerFunc(pattern string, fn QueueFunc) func(ctx context.Context, task *asynq.Task) error {
	return func(ctx context.Context, task *asynq.Task) error {
		err := fn(ctx, task)
		if err != nil {
			log.Println("Error executing handler func for pattern", pattern, err)
			return err
		}

		log.Println("Successfully executed handler func for pattern", pattern)
		return nil
	}
}

// Run starts the asynq server with the configured ServerMux.
func (as *AsynqServer) Run() {
	go func() {
		err := as.Server.Run(as.ServerMux)
		if err != nil {
			log.Fatalf("could not start asynq server: %v", err)
		}
	}()

	log.Printf("asynq server is running")
}

// Stop gracefully stops the asynq server.
func (as *AsynqServer) Stop() {
	// Create a channel to listen for termination signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create a context for the server shutdown
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start a goroutine to listen for termination signals
	go func() {
		<-sigChan // Wait for a signal to be received
		cancel()  // Cancel the context to initiate shutdown
	}()

	// Stop the server gracefully using the context
	as.Server.Shutdown()
}
