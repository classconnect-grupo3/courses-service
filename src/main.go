package main

import (
	"courses-service/src/config"
	"courses-service/src/router"
	"fmt"
	"log"
)

func main() {
	config := config.NewConfig()
	r := router.NewRouter(config)
	if err := r.Run(fmt.Sprintf("%s:%s", config.Host, config.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
