package item

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/stretchr/testify/mock"

	"github.com/sborsh1kmusora/micro-shop/inventory/internal/model"
)

func (s *ServiceSuite) TestCreate() {
	tests := []struct {
		name       string
		item       *model.Item
		mock       func()
		expectUUID string
		expectErr  error
	}{
		{
			name: "success create",
			item: &model.Item{Name: "test item"},
			mock: func() {
				s.inventoryRepo.
					On("Create", mock.Anything, mock.AnythingOfType("*model.Item")).
					Return("uuid-123", nil)
			},
			expectUUID: "uuid-123",
			expectErr:  nil,
		},
		{
			name: "repo error",
			item: &model.Item{Name: "bad item"},
			mock: func() {
				s.inventoryRepo.
					On("Create", mock.Anything, mock.AnythingOfType("*model.Item")).
					Return("", errors.New("repo error"))
			},
			expectUUID: "",
			expectErr:  errors.New("repo error"),
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.inventoryRepo.ExpectedCalls = nil
			s.inventoryRepo.Calls = nil

			tt.mock()

			uuid, err := s.service.Create(ctx, tt.item)

			s.Require().Equal(tt.expectUUID, uuid)

			if tt.expectErr != nil {
				s.EqualError(err, tt.expectErr.Error())
			} else {
				s.NoError(err)
			}

			s.inventoryRepo.AssertExpectations(s.T())
		})
	}
}
