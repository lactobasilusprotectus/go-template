package usecase

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/lactobasilusprotectus/go-template/pkg/util/queue"
	"log"
)

func (a *AuthUseCase) RegisterQueue(as *queue.AsynqServer) {
	as.AddHandlerFunc("send_email", a.SendEmail)
}

func (a *AuthUseCase) SendEmail(ctx context.Context, task *asynq.Task) error {
	log.Print("send email")
	return nil
}
