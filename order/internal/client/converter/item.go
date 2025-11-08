package converter

import (
	"github.com/sborsh1kmusora/micro-shop/order/internal/model"
	inventoryV1 "github.com/sborsh1kmusora/micro-shop/shared/pkg/proto/inventory/v1"
)

func ItemsFilterToProto(filter *model.Filter) *inventoryV1.ItemsFilter {
	categories := make([]inventoryV1.Category, 0, len(filter.Categories))
	for _, category := range filter.Categories {
		categories = append(categories, inventoryV1.Category(category))
	}

	return &inventoryV1.ItemsFilter{
		Uuids:                 filter.UUIDs,
		Names:                 filter.Names,
		Categories:            categories,
		ManufacturerCountries: filter.Countries,
		Tags:                  filter.Tags,
	}
}

func ItemListToModel(items []*inventoryV1.Item) []*model.Item {
	res := make([]*model.Item, 0, len(items))
	for _, item := range items {
		res = append(res, itemToModel(item))
	}
	return res
}

func itemToModel(item *inventoryV1.Item) *model.Item {
	return &model.Item{
		UUID:          item.Uuid,
		Name:          item.Name,
		Description:   item.Description,
		Price:         item.Price,
		StockQuantity: item.StockQuantity,
		Category:      model.Category(item.Category),
		Dimensions:    dimensionsRepoToModel(item.Dimensions),
		Manufacturer:  manufacturerRepoToModel(item.Manufacturer),
		Tags:          item.Tags,
		Metadata:      make(map[string]model.Value),
		CreatedAt:     item.CreatedAt.AsTime(),
		UpdatedAt:     item.UpdatedAt.AsTime(),
	}
}

func dimensionsRepoToModel(d *inventoryV1.Dimensions) *model.Dimensions {
	return &model.Dimensions{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
		Weight: d.Weight,
	}
}

func manufacturerRepoToModel(m *inventoryV1.Manufacturer) *model.Manufacturer {
	return &model.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}
