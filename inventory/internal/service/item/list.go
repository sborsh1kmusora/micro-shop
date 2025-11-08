package item

import (
	"context"

	"github.com/sborsh1kmusora/micro-shop/inventory/internal/model"
)

func (s *service) List(ctx context.Context, filter *model.Filter) ([]*model.Item, error) {
	items, err := s.orderRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	if filter == nil || isEmptyFilter(filter) {
		return items, nil
	}

	items = filterByUUID(items, filter.UUIDs)
	items = filterByNames(items, filter.Names)
	items = filterByCategories(items, filter.Categories)
	items = filterByManufacturerCountry(items, filter.Countries)
	items = filterByTags(items, filter.Tags)

	return items, nil
}

func isEmptyFilter(f *model.Filter) bool {
	return len(f.UUIDs) == 0 &&
		len(f.Names) == 0 &&
		len(f.Categories) == 0 &&
		len(f.Countries) == 0 &&
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

func filterByUUID(items []*model.Item, uuids []string) []*model.Item {
	if len(uuids) == 0 {
		return items
	}
	var out []*model.Item
	for _, it := range items {
		if inSlice(it.UUID, uuids) {
			out = append(out, it)
		}
	}
	return out
}

func filterByNames(items []*model.Item, names []string) []*model.Item {
	if len(names) == 0 {
		return items
	}
	var out []*model.Item
	for _, it := range items {
		if inSlice(it.Name, names) {
			out = append(out, it)
		}
	}
	return out
}

func filterByCategories(items []*model.Item, cats []model.Category) []*model.Item {
	if len(cats) == 0 {
		return items
	}
	set := make(map[model.Category]struct{}, len(cats))
	for _, c := range cats {
		set[c] = struct{}{}
	}
	var out []*model.Item
	for _, it := range items {
		if _, ok := set[it.Category]; ok {
			out = append(out, it)
		}
	}
	return out
}

func filterByManufacturerCountry(items []*model.Item, countries []string) []*model.Item {
	if len(countries) == 0 {
		return items
	}
	var out []*model.Item
	for _, it := range items {
		if inSlice(it.Manufacturer.Country, countries) {
			out = append(out, it)
		}
	}
	return out
}

func filterByTags(items []*model.Item, tags []string) []*model.Item {
	if len(tags) == 0 {
		return items
	}
	var out []*model.Item
	for _, it := range items {
		if containsAny(tags, it.Tags) {
			out = append(out, it)
		}
	}
	return out
}
