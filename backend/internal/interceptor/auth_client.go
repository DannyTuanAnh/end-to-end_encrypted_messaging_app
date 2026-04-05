package interceptor

import (
	"context"
	"fmt"

	jwt_app "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/package/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ctxKey string

const (
	CtxUserIDKey ctxKey = "user_id"
	CtxCallerKey ctxKey = "caller"
	CtxAudKey    ctxKey = "aud"
)

func AuthClientInterceptor(privKeyPath string) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if skip[method] {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		userID, _ := ctx.Value(CtxUserIDKey).(int64)
		issService, _ := ctx.Value(CtxCallerKey).(string)
		audService, _ := ctx.Value(CtxAudKey).(string)

		claims := jwt_app.CreatePayload(issService, userID, audService)

		token, err := jwt_app.SignJWT(privKeyPath, claims)
		if err != nil {
			return fmt.Errorf("failed to sign JWT: %v", err)
		}

		md := metadata.New(map[string]string{
			"authorization": "Bearer " + token,
		})

		ctx = metadata.NewOutgoingContext(ctx, md)

		return invoker(ctx, method, req, reply, cc, opts...)
	}

}
