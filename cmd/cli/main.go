package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/cli"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 2. Load environment variables from .env file
	utils.LoadEnv()

	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	cli := cli.NewCLI(db.DB)

	if cli.Run(ctx) {
		return
	}

}
