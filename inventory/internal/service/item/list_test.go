package item

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/stretchr/testify/mock"

	"github.com/sborsh1kmusora/micro-shop/inventory/internal/model"
)

func (s *ServiceSuite) TestList() {
	mockItems := []*model.Item{
		{UUID: "1", Name: "iPhone", Category: model.CategoryElectronics, Manufacturer: &model.Manufacturer{Country: "USA"}, Tags: []string{"phone", "apple"}},
		{UUID: "2", Name: "Galaxy", Category: model.CategoryElectronics, Manufacturer: &model.Manufacturer{Country: "Korea"}, Tags: []string{"phone", "samsung"}},
		{UUID: "3", Name: "MacBook", Category: model.CategoryElectronics, Manufacturer: &model.Manufacturer{Country: "USA"}, Tags: []string{"laptop", "apple"}},
		{UUID: "4", Name: "Shoes", Category: model.CategoryClothing, Manufacturer: &model.Manufacturer{Country: "Italy"}, Tags: []string{"fashion"}},
	}

	tests := []struct {
		name      string
		filter    *model.Filter
		mock      func()
		expect    []*model.Item
		expectErr error
	}{
		{
			name:   "no filter returns all",
			filter: nil,
			mock: func() {
				s.inventoryRepo.On("List", mock.Anything).Return(mockItems, nil)
			},
			expect: mockItems,
		},
		{
			name:   "filter by UUID",
			filter: &model.Filter{UUIDs: []string{"1", "3"}},
			mock: func() {
				s.inventoryRepo.On("List", mock.Anything).Return(mockItems, nil)
			},
			expect: []*model.Item{mockItems[0], mockItems[2]},
		},
		{
			name:   "filter by Name",
			filter: &model.Filter{Names: []string{"Galaxy"}},
			mock: func() {
				s.inventoryRepo.On("List", mock.Anything).Return(mockItems, nil)
			},
			expect: []*model.Item{mockItems[1]},
		},
		{
			name:   "filter by Category",
			filter: &model.Filter{Categories: []model.Category{model.CategoryClothing}},
			mock: func() {
				s.inventoryRepo.On("List", mock.Anything).Return(mockItems, nil)
			},
			expect: []*model.Item{mockItems[3]},
		},
		{
			name:   "filter by Country",
			filter: &model.Filter{Countries: []string{"USA"}},
			mock: func() {
				s.inventoryRepo.On("List", mock.Anything).Return(mockItems, nil)
			},
			expect: []*model.Item{mockItems[0], mockItems[2]},
		},
		{
			name:   "filter by Tags",
			filter: &model.Filter{Tags: []string{"apple"}},
			mock: func() {
				s.inventoryRepo.On("List", mock.Anything).Return(mockItems, nil)
			},
			expect: []*model.Item{mockItems[0], mockItems[2]},
		},
		{
			name:   "repo error",
			filter: &model.Filter{},
			mock: func() {
				s.inventoryRepo.On("List", mock.Anything).Return(nil, errors.New("db error"))
			},
			expectErr: errors.New("db error"),
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.inventoryRepo.ExpectedCalls = nil
			s.inventoryRepo.Calls = nil

			tt.mock()

			items, err := s.service.List(ctx, tt.filter)

			if tt.expectErr != nil {
				s.Error(err)
				s.EqualError(err, tt.expectErr.Error())
				s.Nil(items)
			} else {
				s.NoError(err)
				s.Equal(tt.expect, items)
			}

			s.inventoryRepo.AssertExpectations(s.T())
		})
	}
}
