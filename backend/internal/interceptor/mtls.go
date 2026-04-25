package interceptor

import (
	"context"
	"crypto/x509"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func MTLSIdentityInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		p, ok := peer.FromContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing peer info")
		}

		tlsInfo, ok := p.AuthInfo.(credentials.TLSInfo)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing TLS info")
		}

		if len(tlsInfo.State.VerifiedChains) == 0 || len(tlsInfo.State.VerifiedChains[0]) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing verified client cert chain")
		}

		leaf := tlsInfo.State.VerifiedChains[0][0]
		if leaf == nil {
			return nil, status.Error(codes.Unauthenticated, "missing leaf certificate")
		}

		caller, err := extractServiceIdentity(leaf)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		ctx = context.WithValue(ctx, CtxCallerKey, caller)

		return handler(ctx, req)
	}
}

func extractServiceIdentity(cert *x509.Certificate) (string, error) {
	// Ưu tiên SAN DNS
	if len(cert.DNSNames) > 0 {
		return cert.DNSNames[0], nil
	}

	// fallback CN
	if cert.Subject.CommonName != "" {
		return cert.Subject.CommonName, nil
	}

	return "", fmt.Errorf("missing service identity in certificate")
}
