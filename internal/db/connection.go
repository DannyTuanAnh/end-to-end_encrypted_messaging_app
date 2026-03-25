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
	connStr := config.NewConfigDB().DB_DNS()

	conf, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return fmt.Errorf("error parsing DB config: %w", err)
	}

	conf.MaxConns = 50
	conf.MinConns = 5
	conf.MaxConnLifetime = 30 * time.Minute
	conf.MaxConnIdleTime = 5 * time.Minute
	conf.HealthCheckPeriod = 1 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return fmt.Errorf("error creating DB pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close() // Đóng pool nếu ping thất bại
		return fmt.Errorf("error pinging DB: %w", err)
	}

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
