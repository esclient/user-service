package main

import (
	"context"
	"log"

	"google.golang.org/grpc"

	"github.com/esclient/user-service/internal/userservice/config"
	"github.com/esclient/user-service/internal/userservice/repository"
	"github.com/esclient/user-service/internal/userservice/service"
	"github.com/esclient/user-service/internal/userservice/handler"
)

func main() {
	config := config.LoadConfig()

	ctx := context.Background()

	databaseConn := repository.NewDatabaseConnection(ctx, config.databaseURL)
	repository := repository.NewPostgresUserRepository(databaseConn)
	service := service.NewUserService(repository)
	handler := handler.NewUserHandler(service)

	grpcServer :=

	
}