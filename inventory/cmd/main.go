package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	ctx := context.Background()

	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Error loading .env file: %v\n", err)
		return
	}

	dbURI := os.Getenv("MONGO_URI")

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Printf("Error connecting to mongo db: %v\n", err)
		return
	}
	defer func() {
		err := mongoClient.Disconnect(ctx)
		if err != nil {
			log.Printf("failed to disconnect: %v\n", err)
		}
	}()

	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		log.Printf("failed to ping mongo db: %v\n", err)
		return
	}

	db := mongoClient.Database("inventory")

	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.LoggerInterceptor()),
	)

	repo := inventoryRepo.NewRepository(db)
	service := inventoryService.NewService(repo)
	api := inventoryApi.NewApi(service)

	inventoryV1.RegisterInventoryServiceServer(s, api)

	reflection.Register(s)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}

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
