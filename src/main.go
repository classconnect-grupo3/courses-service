package main

import (
	"courses-service/src/config"
	"courses-service/src/router"
	"fmt"
	"log"

	_ "courses-service/src/docs"

	_ "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger"
)

// @title Courses Service API
// @version 1.0
// @description API for managing courses and related resources

// @contact.name   El mejor grupo de todos ndea deau
// @contact.url    https://github.com/classconnect-grupo3
// @contact.email  classconnectingsoft2@gmail.com

func main() {
	config := config.NewConfig()
	r := router.NewRouter(config)
	if err := r.Run(fmt.Sprintf("%s:%s", config.Host, config.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
