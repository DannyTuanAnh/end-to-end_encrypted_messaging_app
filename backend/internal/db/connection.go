package db

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/cloudsqlconn"
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
	dsn := connDB.DB_DNS()

	conf, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("error parsing DB config: %w", err)
	}

	conf.MaxConns = 50
	conf.MinConns = 5
	conf.MaxConnLifetime = 30 * time.Minute
	conf.MaxConnIdleTime = 5 * time.Minute
	conf.HealthCheckPeriod = 1 * time.Minute

	// 2. Nếu là môi trường Cloud (Host chứa dấu ":")
	if strings.Contains(connDB.DB.Host, ":") {
		log.Printf("Using Cloud SQL Connector for: %s", connDB.DB.Host)
		d, err := cloudsqlconn.NewDialer(context.Background())
		if err != nil {
			return err
		}
		// Ép thư viện dùng bộ quay số tự động của Google
		conf.ConnConfig.DialFunc = func(ctx context.Context, _, _ string) (net.Conn, error) {
			return d.Dial(ctx, connDB.DB.Host)
		}
	} else {
		log.Println("Using standard TCP/IP connection for database")
		// Chạy Local thì gán Host/Port bình thường
		conf.ConnConfig.Host = connDB.DB.Host

		p, err := strconv.ParseUint(connDB.DB.Port, 10, 16)
		if err != nil {
			return fmt.Errorf("invalid db port %q: %w", connDB.DB.Port, err)
		}
		conf.ConnConfig.Port = uint16(p)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return fmt.Errorf("error creating DB pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close() // Đóng pool nếu ping thất bại
		return fmt.Errorf("error pinging DB: %w", err)
	}

	log.Println("Database connection established successfully. Ping successful.")

	// Gán pool và khởi tạo sqlc.Queries
	DBPool = pool
	DB = sqlc.New(DBPool)

	log.Println("Connecting to database successfully")

	return nil
}

// Close đóng connection pool (gọi khi shutdown app)
func Close() {
	if DBPool != nil {
		DBPool.Close()
		log.Println("Database connection closed")
	}
}
