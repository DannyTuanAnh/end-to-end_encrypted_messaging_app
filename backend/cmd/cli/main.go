package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/cli"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db"
	redis_memory "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/redis"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 1. Load environment variables from .env file
	utils.LoadEnv()

	// 2. Initialize database connection
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
		return
	}
	defer db.Close()

	// 3. Initialize Redis connection
	rdb, err := redis_memory.InitRedis()
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
		return
	}
	defer rdb.CloseRedis()

	cli := cli.NewCLI(db.DB, rdb.RDB)

	if cli.Run(ctx) {
		return
	}

}
