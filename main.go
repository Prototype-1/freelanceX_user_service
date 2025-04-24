package main

import (
	"log"
	"net"

	"github.com/Prototype-1/freelanceX_user_service/internal/auth/handler"
	"github.com/Prototype-1/freelanceX_user_service/internal/auth/repository"
	"github.com/Prototype-1/freelanceX_user_service/internal/auth/service"
	"github.com/Prototype-1/freelanceX_user_service/pkg/db"
	"github.com/Prototype-1/freelanceX_user_service/pkg/redis"
	"github.com/Prototype-1/freelanceX_user_service/config"
	authPb "github.com/Prototype-1/freelanceX_user_service/proto/auth"

	"google.golang.org/grpc"
)

func main() {
	config.LoadConfig()
	redis.InitRedis()
	dsn := buildDSN()
	dbConn, err := db.InitDB(dsn)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	userRepo := repository.NewUserRepository(dbConn)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	grpcServer := grpc.NewServer()
	authPb.RegisterAuthServiceServer(grpcServer, authHandler)

	listener, err := net.Listen("tcp", ":"+config.AppConfig.Port)
	if err != nil {
		log.Fatalf("failed to listen on port %v: %v", config.AppConfig.Port, err)
	}

	log.Printf("gRPC server listening on :%v", config.AppConfig.Port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}

func buildDSN() string {
	return "host=" + config.AppConfig.DBHost +
		" user=" + config.AppConfig.DBUser +
		" password=" + config.AppConfig.DBPassword +
		" dbname=" + config.AppConfig.DBName +
		" port=" + config.AppConfig.DBPort +
		" sslmode=disable"
}
