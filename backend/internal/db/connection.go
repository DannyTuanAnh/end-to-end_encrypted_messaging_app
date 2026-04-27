package db

import (
	"context"
	"fmt"
	"log"
	"net"
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

// deploy
func InitDB() error {
	ctx := context.Background()
	connDB := config.NewConfigDB()

	// Cloud SQL Connector chỉ cần user/password/dbname
	// KHÔNG dùng host=10.54.80.3 nữa
	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=disable",
		connDB.DB.User,
		connDB.DB.Password,
		connDB.DB.DBName,
	)

	conf, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("error parsing DB config: %w", err)
	}

	// Pool config
	conf.MaxConns = 20
	conf.MinConns = 5
	conf.MaxConnLifetime = 30 * time.Minute
	conf.MaxConnIdleTime = 10 * time.Minute
	conf.HealthCheckPeriod = 1 * time.Minute

	// connDB.DB.Host bây giờ phải là:
	// project:region:instance
	// ví dụ: chat-app-493208:us-central1:chat-app-db
	instanceConnName := connDB.DB.Host

	log.Printf("Using Cloud SQL Connector for instance: %s", instanceConnName)

	// Tạo Cloud SQL Dialer
	dialer, err := cloudsqlconn.NewDialer(
		ctx,
		cloudsqlconn.WithLazyRefresh(),
		cloudsqlconn.WithDefaultDialOptions(
			cloudsqlconn.WithPrivateIP(),
		),
	)
	if err != nil {
		return fmt.Errorf("error creating Cloud SQL dialer: %w", err)
	}

	// Override pgx dialer:
	// thay vì dial host:port bình thường
	// pgx sẽ dùng Cloud SQL connector
	conf.ConnConfig.DialFunc = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dialer.Dial(ctx, instanceConnName)
	}

	connectCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(connectCtx, conf)
	if err != nil {
		return fmt.Errorf("error creating DB pool: %w", err)
	}

	if err := pool.Ping(connectCtx); err != nil {
		pool.Close()
		return fmt.Errorf("error pinging DB: %w", err)
	}

	log.Println("DATABASE CONNECTED SUCCESSFULLY VIA CLOUD SQL CONNECTOR")

	DBPool = pool
	DB = sqlc.New(DBPool)

	log.Println("Database initialized successfully")
	return nil
}

func isValidIP(host string) bool {
	return net.ParseIP(host) != nil
}

//local test
// func InitDB() error {
// 	connStr := config.NewConfigDB().DB_DNS()

// 	conf, err := pgxpool.ParseConfig(connStr)
// 	if err != nil {
// 		return fmt.Errorf("error parsing DB config: %w", err)
// 	}

// 	conf.MaxConns = 50
// 	conf.MinConns = 5
// 	conf.MaxConnLifetime = 30 * time.Minute
// 	conf.MaxConnIdleTime = 5 * time.Minute
// 	conf.HealthCheckPeriod = 1 * time.Minute

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	pool, err := pgxpool.NewWithConfig(ctx, conf)
// 	if err != nil {
// 		return fmt.Errorf("error creating DB pool: %w", err)
// 	}

// 	if err := pool.Ping(ctx); err != nil {
// 		pool.Close() // Đóng pool nếu ping thất bại
// 		return fmt.Errorf("error pinging DB: %w", err)
// 	}

// 	// Gán pool và khởi tạo sqlc.Queries
// 	DBPool = pool
// 	DB = sqlc.New(DBPool)

// 	log.Println("Connecting to database successfully")

// 	return nil
// }

// Close đóng connection pool (gọi khi shutdown app)
func Close() {
	if DBPool != nil {
		DBPool.Close()
		log.Println("Database connection closed")
	}
}
