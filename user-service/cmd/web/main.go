package main

import (
	"fmt"
	"log"

	"github.com/DevisArya/learn-microservices/user-service/internal/config"

	"github.com/go-playground/validator/v10"
)

func main() {

	validate := validator.New()
	db := config.NewDB()

	bootstrapResult, err := config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		Validate: validate,
	})

	if err != nil {
		log.Fatalf("failed to bootstrap: %v", err)
	}

	fmt.Println("gRPC server running on port 50051")

	if err := bootstrapResult.GRPCServer.Serve(bootstrapResult.Listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
