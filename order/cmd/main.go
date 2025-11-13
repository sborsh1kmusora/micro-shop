package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderApi "github.com/sborsh1kmusora/micro-shop/order/internal/api/order/v1"
	inventoryClient "github.com/sborsh1kmusora/micro-shop/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/sborsh1kmusora/micro-shop/order/internal/client/grpc/payment/v1"
	customMiddleware "github.com/sborsh1kmusora/micro-shop/order/internal/middleware"
	"github.com/sborsh1kmusora/micro-shop/order/internal/migrator"
	orderRepo "github.com/sborsh1kmusora/micro-shop/order/internal/repository/order"
	orderService "github.com/sborsh1kmusora/micro-shop/order/internal/service/order"
	orderV1 "github.com/sborsh1kmusora/micro-shop/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/sborsh1kmusora/micro-shop/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/sborsh1kmusora/micro-shop/shared/pkg/proto/payment/v1"
)

const (
	httpPort             = "8080"
	inventoryServiceAddr = "localhost:50051"
	paymentServiceAddr   = "localhost:50052"
	readHeaderTimeout    = time.Second * 10
	shutdownTimeout      = time.Second * 10
)

func main() {
	ctx := context.Background()

	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("failed to load .env file: %v\n", err)
		return
	}

	dbURI := os.Getenv("DB_URI")

	pool, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		log.Printf("failed to connect to database: %v\n", err)
		return
	}
	defer pool.Close()

	if err = pool.Ping(ctx); err != nil {
		log.Printf("failed to ping database: %v\n", err)
		return
	}

	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	migratorRunner := migrator.NewMigrator(stdlib.OpenDB(*pool.Config().ConnConfig), migrationsDir)

	err = migratorRunner.Up()
	if err != nil {
		log.Printf("database migration error: %v\n", err)
		return
	}

	invConn, err := grpc.NewClient(
		inventoryServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to create gRPC connection to inventory service: %v", err)
		return
	}

	paymentConn, err := grpc.NewClient(
		paymentServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to create gRPC connection to payment service: %v", err)
		return
	}

	invGrpc := inventoryV1.NewInventoryServiceClient(invConn)
	payGrpc := paymentV1.NewPaymentServiceClient(paymentConn)

	payment := paymentClient.NewClient(payGrpc)
	inventory := inventoryClient.NewClient(invGrpc)

	repo := orderRepo.NewOrderRepository(pool)
	service := orderService.NewService(repo, inventory, payment)

	api := orderApi.NewApi(service)

	orderServer, err := orderV1.NewServer(api)
	if err != nil {
		log.Printf("Error creating order server: %v", err)
		return
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(customMiddleware.RequestLogger)

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:        net.JoinHostPort("localhost", httpPort),
		Handler:     r,
		ReadTimeout: readHeaderTimeout,
	}

	go func() {
		log.Println("Starting server on", httpPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error starting HTTP server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down server: %v", err)
	}

	log.Println("Server gracefully stopped")
}
