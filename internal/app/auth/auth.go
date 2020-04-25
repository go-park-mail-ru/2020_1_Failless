package auth

import (
	"failless/internal/pkg/logger"
	"failless/internal/pkg/settings"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"

	pb "failless/api/proto/auth"
	conf "failless/configs/auth"
	auth "failless/internal/pkg/auth/usecase"
)

func Start() {
	file := logger.OpenLogFile("auth")
	defer file.Close()

	if ok := settings.CheckSecretes(conf.Secrets); !ok {
		log.Println("Can't find variables ", conf.Secrets)
		log.Fatal("Environment variables don't set")
	}
	serverSettings := conf.GetConfig()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", serverSettings.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	srv := auth.GetUseCase()
	pb.RegisterAuthServer(grpcServer, &srv)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop
}
