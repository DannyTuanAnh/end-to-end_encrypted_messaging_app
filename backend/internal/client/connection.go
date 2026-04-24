package client

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/interceptor"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func NewGRPCConn(addr, serverName, certFile, keyFile, keyClient string) (*grpc.ClientConn, error) {
	certPEM := []byte(certFile)
	keyPEM := []byte(keyFile)

	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, err
	}

	caCert := []byte(utils.GetEnv("PATH_CERT_CA", ""))

	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caPool,

		ServerName: serverName,
	}

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
		grpc.WithUnaryInterceptor(interceptor.AuthClientInterceptor(keyClient)),
	)

	if err != nil {
		return nil, err
	}

	return conn, nil
}
