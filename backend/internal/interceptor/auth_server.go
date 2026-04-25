package interceptor

import (
	"context"
	"strings"

	jwt_app "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/package/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var skip = map[string]bool{
	"/proto.AuthService/LoginGoogle": true,
	"/proto.AuthService/Logout":      true,
}

func JWTAuthServerInterceptor(certPath string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		if skip[info.FullMethod] {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		auth := md["authorization"]
		if len(auth) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing token")
		}

		tokenString := strings.TrimPrefix(auth[0], "Bearer ")

		token, err := jwt_app.VerifyJWT(certPath, tokenString)
		if err != nil || !token.Valid {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		claims, ok := token.Claims.(*jwt_app.CustomClaims)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "invalid token claims")
		}

		ctx = context.WithValue(ctx, CtxUserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, CtxCallerKey, claims.Issuer)
		ctx = context.WithValue(ctx, CtxAudKey, claims.Audience)

		return handler(ctx, req)
	}
}
