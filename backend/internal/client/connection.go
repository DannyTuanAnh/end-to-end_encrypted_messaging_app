package client

import (
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
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

	log.Println("Loaded CA cert for gRPC client")
	log.Println(caCertStr)

	caCert := []byte(caCertStr)

	caPool := x509.NewCertPool()
	// caPool.AppendCertsFromPEM(caCert)

	if ok := caPool.AppendCertsFromPEM(caCert); !ok {
		log.Println("ERROR: Could not append CA certs. Check PATH_CERT_CA format.")
		// Nếu thất bại, có thể do thiếu dấu xuống dòng, hãy thử log ra để check
		return nil, fmt.Errorf("failed to append CA certificates")
	}

	block, _ := pem.Decode(caCert)
	if block == nil {
		log.Fatal("failed to parse CA PEM")
	}

	certParsed, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Gateway trusted CA Subject=%s", certParsed.Subject)
	log.Printf("Gateway trusted CA Issuer=%s", certParsed.Issuer)
	log.Printf("Gateway trusted CA Serial=%s", certParsed.SerialNumber.String())
	log.Printf("Gateway trusted CA IsCA=%v", certParsed.IsCA)
	log.Printf("Gateway trusted CA SHA256=%x", sha256.Sum256(certParsed.Raw))

	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		host = addr
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caPool,

		ServerName: serverName, // Quan trọng: phải khớp với CN/SAN của server certificate

		MinVersion: tls.VersionTLS12,
	}

	// Tách lấy host từ addr (bỏ port)
	log.Printf("DEBUG: Attempting to connect to gRPC server at %s with TLS. ServerName: %s", addr, serverName)

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
		grpc.WithUnaryInterceptor(interceptor.AuthClientInterceptor(keyClient)),
		grpc.WithAuthority(host), // route Cloud Run
	)

	if err != nil {
		return nil, err
	}

	return conn, nil
}
