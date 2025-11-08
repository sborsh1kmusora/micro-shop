package converter

import (
	"time"

	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/sborsh1kmusora/micro-shop/inventory/internal/model"
	inventoryV1 "github.com/sborsh1kmusora/micro-shop/shared/pkg/proto/inventory/v1"
)

func ItemToProto(item *model.Item) *inventoryV1.Item {
	var updatedAt *timestamppb.Timestamp
	if item.UpdatedAt != nil {
		updatedAt = timestamppb.New(*item.UpdatedAt)
	}

	return &inventoryV1.Item{
		Uuid:          item.UUID,
		Name:          item.Name,
		Description:   item.Description,
		Price:         item.Price,
		StockQuantity: item.StockQuantity,
		Category:      inventoryV1.Category(item.Category),
		Dimensions:    dimensionsToProto(item.Dimensions),
		Manufacturer:  manufacturerToProto(item.Manufacturer),
		Tags:          item.Tags,
		Metadata:      make(map[string]*inventoryV1.Value),
		CreatedAt:     timestamppb.New(item.CreatedAt),
		UpdatedAt:     updatedAt,
	}
}

func ItemProtoToModel(item *inventoryV1.Item) *model.Item {
	var updatedAt *time.Time
	if item.UpdatedAt != nil {
		updatedAt = lo.ToPtr(item.UpdatedAt.AsTime())
	}

	return &model.Item{
		UUID:          item.Uuid,
		Name:          item.Name,
		Description:   item.Description,
		Price:         item.Price,
		StockQuantity: item.StockQuantity,
		Category:      model.Category(item.Category),
		Dimensions:    dimensionsProtoToModel(item.Dimensions),
		Manufacturer:  manufacturerProtoToModel(item.Manufacturer),
		Tags:          item.Tags,
		Metadata:      make(map[string]*model.Value),
		CreatedAt:     item.CreatedAt.AsTime(),
		UpdatedAt:     updatedAt,
	}
}

func FilterProtoToModel(f *inventoryV1.ItemsFilter) *model.Filter {
	return &model.Filter{
		UUIDs:      f.Uuids,
		Names:      f.Names,
		Tags:       f.Tags,
		Categories: convertCategories(f.Categories),
		Countries:  f.ManufacturerCountries,
	}
}

func ListItemsToProto(items []*model.Item) []*inventoryV1.Item {
	res := make([]*inventoryV1.Item, len(items))
	for i, item := range items {
		res[i] = ItemToProto(item)
	}
	return res
}

func convertCategories(c []inventoryV1.Category) []model.Category {
	res := make([]model.Category, len(c))
	for i, v := range c {
		res[i] = model.Category(v) // cast enum → enum
	}
	return res
}

func dimensionsToProto(d *model.Dimensions) *inventoryV1.Dimensions {
	return &inventoryV1.Dimensions{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
		Weight: d.Weight,
	}
}

func manufacturerToProto(m *model.Manufacturer) *inventoryV1.Manufacturer {
	return &inventoryV1.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

func dimensionsProtoToModel(d *inventoryV1.Dimensions) *model.Dimensions {
	return &model.Dimensions{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
		Weight: d.Weight,
	}
}

func manufacturerProtoToModel(m *inventoryV1.Manufacturer) *model.Manufacturer {
	return &model.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}
