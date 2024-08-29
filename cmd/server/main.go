package main

import (
	"fmt"

	"github.com/RhinoSC/sre-backend/internal/application"
)

func main() {
	port := 8080
	portS := fmt.Sprintf(":%d", port)
	cfg := application.ConfigServerChi{
		Address: portS,
	}
	server := application.NewServerChi(cfg)

	fmt.Printf("listening to: %d\n", port)
	if err := server.Run(); err != nil {
		fmt.Println(err)
		return
	}
}
