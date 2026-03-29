package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/app"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db"
	redis_memory "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/redis"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
)

func main() {
	// Initialize original context for the application
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

	// 4. Initialize application
	application := app.NewApplication(ctx, db.DB, rdb.RDB)

	// 5. Run the application and capture any error messages
	// Check if the application is running in development mode by checking the ENV environment variable
	if ENV := os.Getenv("ENV"); ENV == "development" {
		// RunTLS the application and capture any error messages
		msg, err := application.RunTLS(ctx)
		if err != nil {
			log.Fatalf("%s: %v\n", msg, err)
		}

		log.Println(msg)
	} else {
		// Run the application and capture any error messages
		msg, err := application.Run(ctx)
		if err != nil {
			log.Fatalf("%s: %v\n", msg, err)
		}

		log.Println(msg)
	}

}
