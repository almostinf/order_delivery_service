package main

import (
	"log"

	"github.com/almostinf/order_delivery_service/config"
	"github.com/almostinf/order_delivery_service/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
