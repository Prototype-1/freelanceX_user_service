// Initial setup for GRPC-based User Service

package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	"github.com/Prototype-1/freelanceX_user_service/pkg/db"
	"github.com/Prototype-1/freelanceX_user_service/pkg/logger"
	"github.com/Prototype-1/freelanceX_user_service/pkg/redis"
	
	authPb "github.com/Prototype-1/freelanceX_user_service/proto/auth"
	authHandler "github.com/Prototype-1/freelanceX_user_service/internal/auth/handler"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	db.InitDB()
	redis.InitRedis()
	logger.InitLogger()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	
	grpcServer := grpc.NewServer()

	// Register GRPC auth service
	authPb.RegisterAuthServiceServer(grpcServer, &authHandler.AuthService{})
	log.Printf("gRPC server listening on port %s", port)
	
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}