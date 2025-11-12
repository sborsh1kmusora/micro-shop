package item

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/sborsh1kmusora/micro-shop/inventory/internal/model"
)

func (s *ServiceSuite) TestGet() {
	tests := []struct {
		name       string
		uuid       string
		mock       func()
		expectItem *model.Item
		expectErr  error
	}{
		{
			name: "success get item",
			uuid: "uuid-123",
			mock: func() {
				s.inventoryRepo.
					On("Get", mock.Anything, "uuid-123").
					Return(&model.Item{
						UUID: "uuid-123",
						Name: "test item",
					}, nil)
			},
			expectItem: &model.Item{UUID: "uuid-123", Name: "test item"},
			expectErr:  nil,
		},
		{
			name: "item not found",
			uuid: "uuid-404",
			mock: func() {
				s.inventoryRepo.
					On("Get", mock.Anything, "uuid-404").
					Return((*model.Item)(nil), model.ErrItemNotFound)
			},
			expectItem: nil,
			expectErr:  model.ErrItemNotFound,
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.inventoryRepo.ExpectedCalls = nil
			s.inventoryRepo.Calls = nil

			tt.mock()

			item, err := s.service.Get(ctx, tt.uuid)

			if tt.expectErr != nil {
				s.Error(err)
				s.EqualError(err, tt.expectErr.Error())
				s.Nil(item)
			} else {
				s.NoError(err)
				s.Equal(tt.expectItem, item)
			}

			s.inventoryRepo.AssertExpectations(s.T())
		})
	}
}
