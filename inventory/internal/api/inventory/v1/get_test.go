package v1

import (
	"context"
	"time"

	"github.com/go-faster/errors"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/sborsh1kmusora/micro-shop/inventory/internal/model"
	inventoryV1 "github.com/sborsh1kmusora/micro-shop/shared/pkg/proto/inventory/v1"
)

func (s *APISuite) TestGetItem() {
	tests := []struct {
		name      string
		req       *inventoryV1.GetItemRequest
		mock      func()
		expect    *inventoryV1.GetItemResponse
		expectErr error
	}{
		{
			name: "success",
			req:  &inventoryV1.GetItemRequest{Uuid: "123"},
			mock: func() {
				s.inventoryService.
					On("Get", mock.Anything, "123").
					Return(&model.Item{
						UUID:         "123",
						Name:         "iPhone",
						Category:     model.CategoryElectronics,
						Manufacturer: &model.Manufacturer{Country: "USA"},
						Tags:         []string{"apple"},
						Dimensions:   &model.Dimensions{},
					}, nil)
			},
			expect: &inventoryV1.GetItemResponse{
				Item: &inventoryV1.Item{
					Uuid:     "123",
					Name:     "iPhone",
					Category: inventoryV1.Category_CATEGORY_ELECTRONICS,
					Manufacturer: &inventoryV1.Manufacturer{
						Country: "USA",
					},
					Tags:       []string{"apple"},
					Dimensions: &inventoryV1.Dimensions{},
					Metadata:   map[string]*inventoryV1.Value{},
					CreatedAt:  timestamppb.New(time.Time{}),
				},
			},
			expectErr: nil,
		},
		{
			name: "item not found -> NotFound code",
			req:  &inventoryV1.GetItemRequest{Uuid: "404"},
			mock: func() {
				s.inventoryService.
					On("Get", mock.Anything, "404").
					Return(nil, model.ErrItemNotFound)
			},
			expect:    nil,
			expectErr: status.Error(codes.NotFound, "item not found"),
		},
		{
			name: "internal service error",
			req:  &inventoryV1.GetItemRequest{Uuid: "err"},
			mock: func() {
				s.inventoryService.
					On("Get", mock.Anything, "err").
					Return(nil, errors.New("db broken"))
			},
			expect:    nil,
			expectErr: errors.New("db broken"),
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		s.Run(tt.name, func() {
			// reset mock
			s.inventoryService.ExpectedCalls = nil
			s.inventoryService.Calls = nil

			tt.mock()

			resp, err := s.api.GetItem(ctx, tt.req)

			if tt.expectErr != nil {
				s.Error(err)

				// gRPC errors compare via status.Code
				s.Equal(status.Code(tt.expectErr), status.Code(err))
				s.Equal(tt.expectErr.Error(), err.Error())
			} else {
				s.NoError(err)
				s.Equal(tt.expect, resp)
			}

			s.inventoryService.AssertExpectations(s.T())
		})
	}
}
