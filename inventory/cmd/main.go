package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/sborsh1kmusora/micro-shop/inventory/internal/interceptor"
	inventoryV1 "github.com/sborsh1kmusora/micro-shop/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50051

type inventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer

	mu    sync.RWMutex
	items map[string]*inventoryV1.Item
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.LoggerInterceptor()),
	)

	invService := &inventoryService{
		items: make(map[string]*inventoryV1.Item),
	}

	inventoryV1.RegisterInventoryServiceServer(s, invService)

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

func (s *inventoryService) GetItem(
	ctx context.Context,
	req *inventoryV1.GetItemRequest,
) (*inventoryV1.GetItemResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.items[req.GetUuid()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "item with UUID %s not found", req.GetUuid())
	}

	return &inventoryV1.GetItemResponse{
		Item: item,
	}, nil
}

func (s *inventoryService) ListItems(
	ctx context.Context,
	req *inventoryV1.ListItemsRequest,
) (*inventoryV1.ListItemsResponse, error) {
	filter := req.GetFilter()

	s.mu.RLock()
	items := make([]*inventoryV1.Item, 0, len(s.items))
	for _, it := range s.items {
		items = append(items, it)
	}
	s.mu.RUnlock()

	if filter == nil || isEmptyFilter(filter) {
		return &inventoryV1.ListItemsResponse{Items: items}, nil
	}

	// Последовательно применяем фильтры (AND между полями)
	items = filterByUUID(items, filter.Uuids)
	items = filterByNames(items, filter.Names)
	items = filterByCategories(items, filter.Categories)
	items = filterByManufacturerCountry(items, filter.ManufacturerCountries)
	items = filterByTags(items, filter.Tags)

	return &inventoryV1.ListItemsResponse{Items: items}, nil
}

func (s *inventoryService) AddItem(
	ctx context.Context,
	req *inventoryV1.AddItemRequest,
) (*inventoryV1.AddItemResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	newUUID := uuid.NewString()

	item := &inventoryV1.Item{
		Uuid:          newUUID,
		Name:          req.GetItem().GetName(),
		Description:   req.GetItem().GetDescription(),
		Price:         req.GetItem().GetPrice(),
		StockQuantity: req.GetItem().GetStockQuantity(),
		Category:      req.GetItem().GetCategory(),
		Dimensions:    req.GetItem().GetDimensions(),
		Manufacturer:  req.GetItem().GetManufacturer(),
		Tags:          req.GetItem().GetTags(),
		Metadata:      req.GetItem().GetMetadata(),
		CreatedAt:     timestamppb.New(time.Now()),
	}

	s.items[newUUID] = item

	return &inventoryV1.AddItemResponse{Uuid: newUUID}, nil
}

func isEmptyFilter(f *inventoryV1.ItemsFilter) bool {
	return len(f.Uuids) == 0 &&
		len(f.Names) == 0 &&
		len(f.Categories) == 0 &&
		len(f.ManufacturerCountries) == 0 &&
		len(f.Tags) == 0
}

func inSlice[T comparable](target T, arr []T) bool {
	for _, v := range arr {
		if v == target {
			return true
		}
	}
	return false
}

func containsAny(needle, haystack []string) bool {
	for _, v := range haystack {
		if inSlice(v, needle) {
			return true
		}
	}
	return false
}

func filterByUUID(items []*inventoryV1.Item, uuids []string) []*inventoryV1.Item {
	if len(uuids) == 0 {
		return items
	}
	var out []*inventoryV1.Item
	for _, it := range items {
		if inSlice(it.Uuid, uuids) {
			out = append(out, it)
		}
	}
	return out
}

func filterByNames(items []*inventoryV1.Item, names []string) []*inventoryV1.Item {
	if len(names) == 0 {
		return items
	}
	var out []*inventoryV1.Item
	for _, it := range items {
		if inSlice(it.Name, names) {
			out = append(out, it)
		}
	}
	return out
}

func filterByCategories(items []*inventoryV1.Item, cats []inventoryV1.Category) []*inventoryV1.Item {
	if len(cats) == 0 {
		return items
	}
	set := make(map[inventoryV1.Category]struct{}, len(cats))
	for _, c := range cats {
		set[c] = struct{}{}
	}
	var out []*inventoryV1.Item
	for _, it := range items {
		if _, ok := set[it.Category]; ok {
			out = append(out, it)
		}
	}
	return out
}

func filterByManufacturerCountry(items []*inventoryV1.Item, countries []string) []*inventoryV1.Item {
	if len(countries) == 0 {
		return items
	}
	var out []*inventoryV1.Item
	for _, it := range items {
		if inSlice(it.Manufacturer.GetCountry(), countries) {
			out = append(out, it)
		}
	}
	return out
}

func filterByTags(items []*inventoryV1.Item, tags []string) []*inventoryV1.Item {
	if len(tags) == 0 {
		return items
	}
	var out []*inventoryV1.Item
	for _, it := range items {
		if containsAny(tags, it.Tags) {
			out = append(out, it)
		}
	}
	return out
}
