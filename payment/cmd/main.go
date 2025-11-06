package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/sborsh1kmusora/micro-shop/payment/internal/interceptor"
	paymentV1 "github.com/sborsh1kmusora/micro-shop/shared/pkg/proto/payment/v1"
)

const grpcPort = 50052

type paymentService struct {
	paymentV1.UnimplementedPaymentServiceServer
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.LoggerInterceptor()),
	)

	paymentV1.RegisterPaymentServiceServer(grpcServer, &paymentService{})

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

func (s *paymentService) PayOrder(
	ctx context.Context,
	req *paymentV1.PayOrderRequest,
) (*paymentV1.PayOrderResponse, error) {
	transactionUUID := uuid.New()

	log.Printf("Оплата прошла успешно, transacion uuid: %s", transactionUUID.String())

	return &paymentV1.PayOrderResponse{
		TransactionUuid: transactionUUID.String(),
	}, nil
}
