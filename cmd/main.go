package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	user_service "github.com/esclient/user-service/api/userservice"

	//Для ветки main убрать ссылку на ветку (/tree/beta)
	"github.com/esclient/user-service/tree/beta/internal/userservice/config"
	"github.com/esclient/user-service/tree/beta/internal/userservice/handler"
	"github.com/esclient/user-service/tree/beta/internal/userservice/repository"
	"github.com/esclient/user-service/tree/beta/internal/userservice/service"
)

func main() {
	cfg := config.LoadConfig()

	ctx := context.Background()

	databaseConn := repository.NewDatabaseConnection(ctx, cfg.DatabaseURL)
	repository := repository.NewPostgresUserRepository(databaseConn)
	
	userService := service.NewUserService(repository)
	userHandler := handler.NewUserHandler(userService)

	grpcServer := grpc.NewServer()

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