package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/app"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/config"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
)

func main() {
	// Initialize original context for the application
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 2. Load environment variables from .env file
	utils.LoadEnv()

	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Check for command-line arguments and execute corresponding commands before starting the application
	hasExecuteCmd := utils.CommandTool(ctx, db.DB)
	if hasExecuteCmd {
		return
	}

	// 1. Initialize configuration
	cfg := config.NewConfigServer()

	// 2. Initialize application
	application := app.NewApplication(ctx, cfg, db.DB)

	// 3. Run the application and capture any error message
	msg, err := application.Run(ctx)
	if err != nil {
		log.Fatalf("%s: %v\n", msg, err)
	}

	log.Println(msg)
}
