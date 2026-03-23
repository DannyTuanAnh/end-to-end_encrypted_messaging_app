package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/app"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/config"
)

func main() {
	// Initialize original context for the application
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 1. Initialize configuration
	cfg := config.NewConfigServer()

	// 2. Initialize application
	application := app.NewApplication(ctx, cfg)

	// 3. Run the application and capture any error message
	msg, err := application.Run(ctx)
	if err != nil {
		log.Fatalf("%s: %v\n", msg, err)
	}

	log.Println(msg)
}
