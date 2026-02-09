package gapi

import (
	"context"
	db "simplebank/db/sqlc"
	"simplebank/pb"
	"simplebank/util"
	worker "simplebank/worker"
	"time"

	"github.com/hibiken/asynq"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	HashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Username:       req.GetUsername(),
			HashedPassword: HashedPassword,
			FullName:       req.GetFullName(),
			Email:          req.GetEmail(),
		},
		AfterCreate: func(user db.User) error {
			payload := &worker.PayloadSendVerifyEmail{Username: user.Username}
			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue("critical"),
			}
			err = server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, payload, opts...)
			if err != nil {
				return status.Error(codes.Internal, "failed to distribute task to send verify email")
			}
			return nil
		},
	}

	result, err := server.store.CreateUserTx(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Error(codes.AlreadyExists, err.Error())
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	rsp := &pb.CreateUserResponse{
		User: convertUser(result.User),
	}
	return rsp, nil
}
