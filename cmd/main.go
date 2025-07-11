package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	user_service "github.com/esclient/user-service/api/userservice"

	"github.com/esclient/user-service/internal/userservice/config"
	"github.com/esclient/user-service/internal/userservice/handler"
	repo "github.com/esclient/user-service/internal/userservice/repository"
	"github.com/esclient/user-service/internal/userservice/service"
)

func main() {
	cfg := config.LoadConfig()

	ctx := context.Background()

	databaseConn, err := repo.NewDatabaseConnection(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	repository := repo.NewPostgresUserRepository(databaseConn)
	
	userService := service.NewUserService(repository)
	userHandler := handler.NewUserHandler(userService)

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	user_service.RegisterUserServiceServer(grpcServer, userHandler)

	listener, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("gRPC server listening at %v", listener.Addr())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}