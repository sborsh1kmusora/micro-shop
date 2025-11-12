package v1

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/stretchr/testify/mock"

	inventoryV1 "github.com/sborsh1kmusora/micro-shop/shared/pkg/proto/inventory/v1"
)

func (s *APISuite) TestAddItem() {
	tests := []struct {
		name      string
		req       *inventoryV1.AddItemRequest
		mock      func()
		expect    *inventoryV1.AddItemResponse
		expectErr error
	}{
		{
			name: "success",
			req: &inventoryV1.AddItemRequest{
				Item: &inventoryV1.Item{
					Uuid:         "123",
					Name:         "iPhone",
					Dimensions:   &inventoryV1.Dimensions{},
					Manufacturer: &inventoryV1.Manufacturer{},
				},
			},
			mock: func() {
				s.inventoryService.
					On("Create", mock.Anything, mock.AnythingOfType("*model.Item")).
					Return("123", nil)
			},
			expect: &inventoryV1.AddItemResponse{Uuid: "123"},
		},
		{
			name: "service error",
			req: &inventoryV1.AddItemRequest{
				Item: &inventoryV1.Item{
					Uuid:         "456",
					Name:         "Samsung",
					Dimensions:   &inventoryV1.Dimensions{},
					Manufacturer: &inventoryV1.Manufacturer{},
				},
			},
			mock: func() {
				s.inventoryService.
					On("Create", mock.Anything, mock.AnythingOfType("*model.Item")).
					Return("", errors.New("db error"))
			},
			expect:    nil,
			expectErr: errors.New("db error"),
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.inventoryService.ExpectedCalls = nil
			s.inventoryService.Calls = nil

			tt.mock()

			resp, err := s.api.AddItem(ctx, tt.req)

			if tt.expectErr != nil {
				s.Error(err)
				s.Equal(tt.expectErr.Error(), err.Error())
				s.Nil(resp)
			} else {
				s.NoError(err)
				s.Equal(tt.expect.Uuid, resp.GetUuid())
			}

			s.inventoryService.AssertExpectations(s.T())
		})
	}
}
