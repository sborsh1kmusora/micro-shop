package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1 "github.com/sborsh1kmusora/micro-shop/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/sborsh1kmusora/micro-shop/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/sborsh1kmusora/micro-shop/shared/pkg/proto/payment/v1"
)

const (
	httpPort             = "8080"
	inventoryServiceAddr = "localhost:50051"
	paymentServiceAddr   = "localhost:50052"
	readHeaderTimeout    = time.Second * 5
	shutdownTimeout      = time.Second * 10
)

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*orderV1.Order
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*orderV1.Order),
	}
}

func (s *OrderStorage) SaveOrder(uuid string, order *orderV1.Order) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[uuid] = order
}

func (s *OrderStorage) LoadOrder(uuid string) (*orderV1.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[uuid]
	if !ok {
		return nil, errors.New("order not found")
	}

	return order, nil
}

type OrderHandler struct {
	inventoryCl inventoryV1.InventoryServiceClient
	paymentCl   paymentV1.PaymentServiceClient
	storage     *OrderStorage
}

func NewOrderHandler(
	storage *OrderStorage,
	invClient inventoryV1.InventoryServiceClient,
	paymentCl paymentV1.PaymentServiceClient,
) *OrderHandler {
	return &OrderHandler{
		storage:     storage,
		inventoryCl: invClient,
		paymentCl:   paymentCl,
	}
}

func (h *OrderHandler) CreateOrder(
	ctx context.Context,
	req *orderV1.CreateOrderRequest,
) (orderV1.CreateOrderRes, error) {
	invResp, err := h.inventoryCl.ListItems(ctx, &inventoryV1.ListItemsRequest{
		Filter: &inventoryV1.ItemsFilter{
			Uuids: req.ItemUuids,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list items: %w", err)
	}

	if len(invResp.Items) != len(req.ItemUuids) {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Один или несколько товаров не найдены",
		}, nil
	}

	var totalPrice float32
	for _, item := range invResp.Items {
		totalPrice += item.Price
	}

	orderUUID, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate uuid: %w", err)
	}

	order := &orderV1.Order{
		OrderUUID:       orderUUID.String(),
		UserUUID:        req.UserUUID,
		ItemUuids:       req.ItemUuids,
		TotalPrice:      totalPrice,
		TransactionUUID: "",
		Status:          orderV1.OrderStatusPendingPayment,
		PaymentMethod:   "",
	}

	h.storage.SaveOrder(orderUUID.String(), order)

	log.Printf("order succeffully created: %+v", order)

	return &orderV1.CreateOrderResponse{
		OrderUUID:  orderUUID.String(),
		TotalPrice: totalPrice,
	}, nil
}

func (h *OrderHandler) GetOrder(
	ctx context.Context,
	params orderV1.GetOrderParams,
) (orderV1.GetOrderRes, error) {
	order, err := h.storage.LoadOrder(params.OrderUUID)
	if err != nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: err.Error(),
		}, nil
	}

	respOrder := orderV1.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		ItemUuids:       order.ItemUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		Status:          order.Status,
		PaymentMethod:   order.PaymentMethod,
	}

	return &orderV1.GetOrderResponse{
		Order: respOrder,
	}, nil
}

func (h *OrderHandler) PayOrder(
	ctx context.Context,
	req *orderV1.PayOrderRequest,
	params orderV1.PayOrderParams,
) (orderV1.PayOrderRes, error) {
	order, err := h.storage.LoadOrder(params.OrderUUID)
	if err != nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: fmt.Sprintf("Заказ с uuid %s не найден", params.OrderUUID),
		}, nil
	}

	if order.Status != orderV1.OrderStatusPendingPayment {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "Заказ уже оплачен или был отменен",
		}, nil
	}

	payResp, err := h.paymentCl.PayOrder(ctx, &paymentV1.PayOrderRequest{
		OrderUuid:     order.OrderUUID,
		UserUuid:      order.UserUUID,
		PaymentMethod: mapPaymentMethod(req.PaymentMethod),
	})
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: fmt.Sprintf("Ошибка платежа: %s", err.Error()),
		}, nil
	}

	order.Status = orderV1.OrderStatusPaid
	order.TransactionUUID = payResp.GetTransactionUuid()
	order.PaymentMethod = req.PaymentMethod

	log.Printf("order succeffully paid: %+v", order)

	return &orderV1.PayOrderResponse{
		TransactionUUID: order.TransactionUUID,
	}, nil
}

func (h *OrderHandler) CancelOrder(
	ctx context.Context,
	params orderV1.CancelOrderParams,
) (orderV1.CancelOrderRes, error) {
	order, err := h.storage.LoadOrder(params.OrderUUID)
	if err != nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: fmt.Sprintf("Заказ с uuid %s не найден", params.OrderUUID),
		}, nil
	}

	if order.Status == orderV1.OrderStatusPaid {
		return &orderV1.ConflictError{
			Code:    409,
			Message: "Заказ уже оплачен и не может быть отменён",
		}, nil
	}

	if order.Status != orderV1.OrderStatusPendingPayment {
		return &orderV1.ConflictError{
			Code:    409,
			Message: "Заказ нельзя отменить",
		}, nil
	}

	order.Status = orderV1.OrderStatusCanceled

	h.storage.SaveOrder(params.OrderUUID, order)

	log.Printf("order succeffully cancelled: %+v", order)

	return &orderV1.CancelOrderNoContent{}, nil
}

func (h *OrderHandler) NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(http.StatusInternalServerError),
			Message: orderV1.NewOptString(err.Error()),
		},
	}
}

func main() {
	storage := NewOrderStorage()

	invConn, err := grpc.NewClient(
		inventoryServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to create gRPC connection to inventory service: %v", err)
	}

	paymentConn, err := grpc.NewClient(
		paymentServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to create gRPC connection to payment service: %v", err)
	}

	inventoryClient := inventoryV1.NewInventoryServiceClient(invConn)
	paymentClient := paymentV1.NewPaymentServiceClient(paymentConn)

	orderHandler := NewOrderHandler(storage, inventoryClient, paymentClient)

	orderServer, err := orderV1.NewServer(orderHandler)
	if err != nil {
		log.Fatalf("Error creating order server: %v", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

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

func mapPaymentMethod(method orderV1.PaymentMethod) paymentV1.PaymentMethod {
	switch method {
	case orderV1.PaymentMethodCard:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case orderV1.PaymentMethodSbp:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_SBP
	case orderV1.PaymentMethodCreditCard:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case orderV1.PaymentMethodInvestorMoney:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}
