package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/lactobasilusprotectus/go-template/pkg/auth/common"
	"github.com/lactobasilusprotectus/go-template/pkg/util/queue"
	"log"
)

func (a *AuthUseCase) RegisterQueue(as *queue.AsynqServer) {
	as.AddHandlerFunc(common.TypeWelcomeEmail, a.HandleSendEmail)
	as.AddHandlerFunc(common.TypeReminderEmail, a.HandleSendEmail)
}

func (a *AuthUseCase) HandleSendEmail(ctx context.Context, task *asynq.Task) error {
	var p common.LoginRequest
	if err := json.Unmarshal(task.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Sending Email to User: user_id=%s, template_id=%s", p.Email, p.Password)
	return nil
}

func (a *AuthUseCase) EmailDeliveryTask(ctx context.Context, request common.LoginRequest) (*asynq.Task, error) {
	payload, err := json.Marshal(request)

	if err != nil {
		return nil, err
	}

	return asynq.NewTask(common.TypeWelcomeEmail, payload), nil
}
