package main

import (
	"courses-service/src/router"
)

func main() {
	r := router.NewRouter()
	r.Run(":8080")
}
