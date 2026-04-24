package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"strings"

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

	caCertStr := utils.GetEnv("PATH_CERT_CA", "")
	// Ép chuỗi "\n" thành dấu xuống dòng thực sự để gRPC đọc được định dạng PEM
	caCertStr = strings.ReplaceAll(caCertStr, "\\n", "\n")

	caCert := []byte(caCertStr)

	caPool := x509.NewCertPool()
	// caPool.AppendCertsFromPEM(caCert)

	if ok := caPool.AppendCertsFromPEM(caCert); !ok {
		log.Println("ERROR: Could not append CA certs. Check PATH_CERT_CA format.")
		// Nếu thất bại, có thể do thiếu dấu xuống dòng, hãy thử log ra để check
		return nil, fmt.Errorf("failed to append CA certificates")
	}

	host, _, _ := net.SplitHostPort(addr)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caPool,

		ServerName:         host,
		InsecureSkipVerify: true,
	}

	// Tách lấy host từ addr (bỏ port)
	log.Printf("DEBUG: Attempting to connect to gRPC server at %s with TLS. ServerName: %s, Host: %s", addr, host, host)

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
		grpc.WithUnaryInterceptor(interceptor.AuthClientInterceptor(keyClient)),
		// QUAN TRỌNG: Ép authority là URL thực tế của Cloud Run
		grpc.WithAuthority(host),
	)

	if err != nil {
		return nil, err
	}

	return conn, nil
}
