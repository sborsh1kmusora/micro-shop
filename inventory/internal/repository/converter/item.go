package converter

import (
	"github.com/sborsh1kmusora/micro-shop/inventory/internal/model"
	repoModel "github.com/sborsh1kmusora/micro-shop/inventory/internal/repository/model"
)

func ItemToRepoModel(item *model.Item) *repoModel.Item {
	return &repoModel.Item{
		UUID:          item.UUID,
		Name:          item.Name,
		Description:   item.Description,
		Price:         item.Price,
		StockQuantity: item.StockQuantity,
		Category:      repoModel.Category(item.Category),
		Dimensions:    dimensionsToRepoModel(item.Dimensions),
		Manufacturer:  manufacturerToRepoModel(item.Manufacturer),
		Tags:          item.Tags,
		Metadata:      make(map[string]*repoModel.Value),
		CreatedAt:     item.CreatedAt,
		UpdatedAt:     item.UpdatedAt,
	}
}

func ItemRepoToModel(item *repoModel.Item) *model.Item {
	return &model.Item{
		UUID:          item.UUID,
		Name:          item.Name,
		Description:   item.Description,
		Price:         item.Price,
		StockQuantity: item.StockQuantity,
		Category:      model.Category(item.Category),
		Dimensions:    dimensionsRepoToModel(item.Dimensions),
		Manufacturer:  manufacturerRepoToModel(item.Manufacturer),
		Tags:          item.Tags,
		Metadata:      make(map[string]*model.Value),
		CreatedAt:     item.CreatedAt,
		UpdatedAt:     item.UpdatedAt,
	}
}

func ItemListToModel(items []*repoModel.Item) []*model.Item {
	res := make([]*model.Item, 0, len(items))
	for _, item := range items {
		res = append(res, ItemRepoToModel(item))
	}
	return res
}

func dimensionsToRepoModel(d *model.Dimensions) *repoModel.Dimensions {
	return &repoModel.Dimensions{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
		Weight: d.Weight,
	}
}

func manufacturerToRepoModel(m *model.Manufacturer) *repoModel.Manufacturer {
	return &repoModel.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

func dimensionsRepoToModel(d *repoModel.Dimensions) *model.Dimensions {
	return &model.Dimensions{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
		Weight: d.Weight,
	}
}

func manufacturerRepoToModel(m *repoModel.Manufacturer) *model.Manufacturer {
	return &model.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}
