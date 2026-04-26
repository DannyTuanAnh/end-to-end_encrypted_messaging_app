package client

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/interceptor"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func NewGRPCConn(addr, serverName, certFile, keyFile, keyClient string) (*grpc.ClientConn, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	caCertStr := utils.GetEnv("PATH_CERT_CA", "")

	log.Println("Loaded CA cert for gRPC client")
	log.Println(caCertStr)

	caCert, err := os.ReadFile(caCertStr)
	if err != nil {
		return nil, err
	}

	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caCert)

	// host, _, err := net.SplitHostPort(addr)
	// if err != nil {
	// 	host = addr
	// }

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caPool,

		ServerName: serverName,

		MinVersion: tls.VersionTLS12,
	}

	// Tách lấy host từ addr (bỏ port)
	log.Printf("DEBUG: Attempting to connect to gRPC server at %s with TLS. ServerName: %s", addr, serverName)

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
		grpc.WithUnaryInterceptor(interceptor.AuthClientInterceptor(keyClient)),
		// grpc.WithAuthority(host), // route Cloud Run
	)

	if err != nil {
		return nil, err
	}

	return conn, nil
}
