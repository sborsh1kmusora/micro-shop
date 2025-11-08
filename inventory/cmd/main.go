package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	inventoryApi "github.com/sborsh1kmusora/micro-shop/inventory/internal/api/inventory/v1"
	"github.com/sborsh1kmusora/micro-shop/inventory/internal/interceptor"
	inventoryRepo "github.com/sborsh1kmusora/micro-shop/inventory/internal/repository/item"
	inventoryService "github.com/sborsh1kmusora/micro-shop/inventory/internal/service/item"
	inventoryV1 "github.com/sborsh1kmusora/micro-shop/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50051

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.LoggerInterceptor()),
	)

	repo := inventoryRepo.NewRepository()
	service := inventoryService.NewService(repo)
	api := inventoryApi.NewApi(service)

	inventoryV1.RegisterInventoryServiceServer(s, api)

	reflection.Register(s)

	go func() {
		log.Printf("starting gRPC server on port %d", grpcPort)
		if err := s.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down gRPC server...")

	s.GracefulStop()

	log.Println("Server gracefully stopped")
}
