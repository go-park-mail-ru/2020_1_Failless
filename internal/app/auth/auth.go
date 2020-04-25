package auth

import (
	"failless/configs/auth"
	"failless/internal/pkg/logger"
	"failless/internal/pkg/settings"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	pb "failless/api/proto/auth"
)

func Start() {
	file := logger.OpenLogFile("auth")
	pb.RegisterAuthServer()
	defer file.Close()

	if ok := settings.CheckSecretes(auth.Secrets); !ok {
		log.Println("Can't find variables ", auth.Secrets)
		log.Fatal("Environment variables don't set")
	}
	serverSettings := auth.GetConfig()

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", serverSettings.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServer(grpcServer, &routeGuideServer{})

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop
}
