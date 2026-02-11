package gapi

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func unauthenticatedError(err error) error {
	return status.Errorf(codes.Unauthenticated, "unauthenticated: %v", err)
}

// 权限不足的错误处理 HTTP 403
func permissionDeniedError(err error) error {
	return status.Errorf(codes.PermissionDenied, "permission denied: %v", err)
}
