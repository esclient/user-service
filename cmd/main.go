package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	user_service "github.com/esclient/user-service/api/userservice"

	"github.com/esclient/user-service/internal/userservice/config"
	"github.com/esclient/user-service/internal/userservice/handler"
	repo "github.com/esclient/user-service/internal/userservice/repository"
	"github.com/esclient/user-service/internal/userservice/service"

	infisical "github.com/esclient/user-service/internal/infisical"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.LoadConfig()

	infisicalClient := infisical.NewClient(infisical.InfisicalURL, cfg.InfisicalSecretKey)

	secretDBURL, err := infisicalClient.GetSecret(ctx, cfg.InfisicalProjectId, cfg.InfisicalEnv, infisical.SecretPathUserService, "DB_URL")
	if err != nil {
		log.Fatal(err)
	}

	databaseConn, err := repo.NewDatabaseConnection(ctx, secretDBURL.Value)
	if err != nil {
		log.Fatal(err)
	}

    repository := repo.NewPostgresUserRepositoryFromPool(databaseConn)
	
	userService := service.NewUserService(repository)
	userHandler := handler.NewUserHandler(userService)

	grpcServer := grpc.NewServer()

	user_service.RegisterUserServiceServer(grpcServer, userHandler)

	reflection.Register(grpcServer)

	address := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("gRPC server listening at %v", listener.Addr())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}