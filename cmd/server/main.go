package main

import (
	"fmt"

	"github.com/RhinoSC/sre-backend/internal/application"
)

func main() {
	fmt.Println("Hello World!")

	cfg := application.ConfigServerChi{
		Address: "localhost:8080",
	}
	server := application.NewServerChi(cfg)

	if err := server.Run(); err != nil {
		fmt.Println(err)
		return
	}
}
