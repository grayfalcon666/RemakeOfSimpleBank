package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/hibiken/asynq"
)

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		slog.ErrorContext(
			ctx,
			"failed to unmarshal verify email task payload",
			"error", err, // 错误详情
			"task_type", task.Type(),
			"task_id", task.ResultWriter().TaskID(),
		)
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	user, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user doesn't exist: %w", asynq.SkipRetry)
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// 模拟发送邮件
	slog.InfoContext(
		ctx,
		"sending verify email to user",
		"username", user.Username,
		"email", user.Email,
		"task_id", task.ResultWriter().TaskID(),
	)
	return nil
}
