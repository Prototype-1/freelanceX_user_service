package main

import (
	"log"
	"net"

	userHdlr "github.com/Prototype-1/freelanceX_user_service/internal/auth/handler"
	userRepo "github.com/Prototype-1/freelanceX_user_service/internal/auth/repository"
	userSvc "github.com/Prototype-1/freelanceX_user_service/internal/auth/service"
	"github.com/Prototype-1/freelanceX_user_service/pkg/db"
	"github.com/Prototype-1/freelanceX_user_service/pkg/redis"
	"github.com/Prototype-1/freelanceX_user_service/config"
	authPb "github.com/Prototype-1/freelanceX_user_service/proto/auth"
	profileHdlr "github.com/Prototype-1/freelanceX_user_service/internal/profile/handler"
   profileRepo "github.com/Prototype-1/freelanceX_user_service/internal/profile/repository"
   profileSvc "github.com/Prototype-1/freelanceX_user_service/internal/profile/service"
profilePb "github.com/Prototype-1/freelanceX_user_service/proto/profile"
portRepo "github.com/Prototype-1/freelanceX_user_service/internal/portfolio/repository"
portSvc "github.com/Prototype-1/freelanceX_user_service/internal/portfolio/service"
portHdlr "github.com/Prototype-1/freelanceX_user_service/internal/portfolio/handler"
portfolioPb "github.com/Prototype-1/freelanceX_user_service/proto/portfolio"
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

	userRepo := userRepo.NewUserRepository(dbConn)
	authService := userSvc.NewAuthService(userRepo)
	authHandler := userHdlr.NewAuthHandler(authService)

	profileRepo := profileRepo.NewRepository(dbConn)
	profileService := profileSvc.NewService(profileRepo)
	profileHandler := profileHdlr.NewHandler(profileService)

	portfolioRepo := portRepo.NewRepository(dbConn)
	portfolioService := portSvc.NewService(portfolioRepo)
	portfolioHandler := portHdlr.NewHandler(portfolioService)

	grpcServer := grpc.NewServer()
	authPb.RegisterAuthServiceServer(grpcServer, authHandler)
	profilePb.RegisterProfileServiceServer(grpcServer, profileHandler)
	portfolioPb.RegisterPortfolioServiceServer(grpcServer, portfolioHandler)


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
