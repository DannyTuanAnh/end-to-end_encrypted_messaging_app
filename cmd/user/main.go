package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/server"
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

	// 5. Initialize application
	userServer := server.NewUserServer(ctx, db.DB)

	// 6. Run the application and capture any error message
	msg, err := userServer.Run()
	if err != nil {
		log.Fatalf("%s: %v\n", msg, err)
	}

	log.Println(msg)
}
