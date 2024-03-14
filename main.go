package main

import (
	"os"

	"github.com/Hy-Iam-Noval/dacol-2/src"
	"github.com/Hy-Iam-Noval/dacol-2/src/ctrl"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r := src.Route()

	ctrl.IsError(r.Listen(":" + port))
}
