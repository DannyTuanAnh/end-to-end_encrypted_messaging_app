package interceptor

import (
	"context"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RBACInterceptor(policies map[string][]string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		caller := utils.GetCaller(ctx)
		if caller == "" || caller == "unknown" {
			return nil, status.Error(codes.PermissionDenied, "caller not identified")
		}

		allowedCallers, ok := policies[info.FullMethod]
		if !ok {
			return nil, status.Error(codes.PermissionDenied, "method not allowed")
		}

		for _, svc := range allowedCallers {
			if svc == caller {
				return handler(ctx, req)
			}
		}

		return nil, status.Errorf(codes.PermissionDenied, "caller %s not allowed to call %s", caller, info.FullMethod)
	}
}
