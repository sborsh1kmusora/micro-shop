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

	paymentApi "github.com/sborsh1kmusora/micro-shop/payment/internal/api/payment/v1"
	"github.com/sborsh1kmusora/micro-shop/payment/internal/interceptor"
	"github.com/sborsh1kmusora/micro-shop/payment/internal/service/payment"
	paymentV1 "github.com/sborsh1kmusora/micro-shop/shared/pkg/proto/payment/v1"
)

const grpcPort = 50052

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.LoggerInterceptor()),
	)

	service := payment.NewPaymentService()
	api := paymentApi.NewApi(service)

	paymentV1.RegisterPaymentServiceServer(grpcServer, api)

	reflection.Register(grpcServer)

	go func() {
		log.Printf("grpc server listening at %d", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down gRPC server...")

	grpcServer.GracefulStop()

	log.Println("Server gracefully stopped")
}
