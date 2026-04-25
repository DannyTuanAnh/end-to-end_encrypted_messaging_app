package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/config"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	DB     *sqlc.Queries
	DBPool *pgxpool.Pool //// Export để có thể Close() khi shutdown
)

func InitDB() error {
	ctx := context.Background()
	connDB := config.NewConfigDB()
	// dsn := connDB.DB_DNS()

	// Tạo DSN đặc biệt cho Unix Socket
	// Lưu ý: host phải là đường dẫn thư mục chứa socket
	instanceName := "chat-app-493208:us-central1:chat-app-db"
	socketPath := "/cloudsql/" + instanceName
	dsn := fmt.Sprintf("user=%s password=%s database=%s host=%s sslmode=disable",
		connDB.DB.User, connDB.DB.Password, connDB.DB.DBName, socketPath)

	log.Println("Connecting via UNIX SOCKET to bypass Mesh Proxy...")

	conf, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("error parsing DB config: %w", err)
	}

	// conf.MaxConns = 50
	// conf.MinConns = 5
	// conf.MaxConnLifetime = 30 * time.Minute
	// conf.MaxConnIdleTime = 5 * time.Minute
	// conf.HealthCheckPeriod = 1 * time.Minute

	// 2. Nếu là môi trường Cloud (Host chứa dấu ":")
	// if strings.Contains(connDB.DB.Host, ":") {
	// 	log.Printf("Using Cloud SQL Connector with PROXY BYPASS for: %s", connDB.DB.Host)

	// 	httpClient := &http.Client{
	// 		Transport: &http.Transport{
	// 			Proxy: nil, // Ép buộc không dùng proxy cho việc lấy metadata
	// 		},
	// 	}

	// 	d, err := cloudsqlconn.NewDialer(
	// 		context.Background(),
	// 		cloudsqlconn.WithDefaultDialOptions(cloudsqlconn.WithPrivateIP()),
	// 		cloudsqlconn.WithHTTPClient(httpClient),
	// 	)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	// Ép thư viện dùng bộ quay số tự động của Google
	// 	conf.ConnConfig.DialFunc = func(ctx context.Context, _, _ string) (net.Conn, error) {
	// 		return d.Dial(ctx, connDB.DB.Host)
	// 	}
	// } else {
	// 	log.Println("Connecting via Private IP with direct TLS...")
	// 	// Chạy Local thì gán Host/Port bình thường
	// 	conf.ConnConfig.Host = connDB.DB.Host

	// 	p, err := strconv.ParseUint(connDB.DB.Port, 10, 16)
	// 	if err != nil {
	// 		return fmt.Errorf("invalid db port %q: %w", connDB.DB.Port, err)
	// 	}
	// 	conf.ConnConfig.Port = uint16(p)

	// 	conf.ConnConfig.TLSConfig = &tls.Config{
	// 		InsecureSkipVerify: true,
	// 	}
	// }

	// // ÉP BUỘC kết nối trực tiếp qua IP (Bỏ qua Cloud SQL Dialer)
	// log.Println("Connecting DIRECTLY to Private IP:", connDB.DB.Host)
	// conf.ConnConfig.Host = connDB.DB.Host
	// p, _ := strconv.ParseUint(connDB.DB.Port, 10, 16)
	// conf.ConnConfig.Port = uint16(p)

	// // BẮT BUỘC có InsecureSkipVerify vì IP 10.54.80.3 không khớp với tên trong Cert của Google
	// conf.ConnConfig.TLSConfig = &tls.Config{
	// 	InsecureSkipVerify: true,
	// }

	connectCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(connectCtx, conf)
	if err != nil {
		return fmt.Errorf("error creating DB pool: %w", err)
	}

	if err := pool.Ping(connectCtx); err != nil {
		pool.Close() // Đóng pool nếu ping thất bại
		return fmt.Errorf("error pinging DB: %w", err)
	}

	log.Println("DATABASE CONNECTED SUCCESSFULLY VIA UNIX SOCKET!")

	// Gán pool và khởi tạo sqlc.Queries
	DBPool = pool
	DB = sqlc.New(DBPool)

	log.Println("Connecting to database successfully")

	return nil
}

// func isValidIP(host string) bool {
// 	return net.ParseIP(host) != nil
// }

// Close đóng connection pool (gọi khi shutdown app)
func Close() {
	if DBPool != nil {
		DBPool.Close()
		log.Println("Database connection closed")
	}
}
