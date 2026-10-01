package inventory

import (
	"context"

	"github.com/DoomGuy1818/rocket-science/inventory/internal/model"
	modelConverter "github.com/DoomGuy1818/rocket-science/inventory/internal/repository/converter"
	repoModel "github.com/DoomGuy1818/rocket-science/inventory/internal/repository/model"
)

func (i *Repository) List(_ context.Context, filters *repoModel.PartsFilter) ([]model.Part, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	response := make([]model.Part, 0, len(i.parts))

	for _, part := range i.parts {
		if filters != nil && !i.matchesFilter(part, filters) {
			continue
		}

		response = append(response, modelConverter.PartToModel(*part))
	}

	return response, nil
}

func (i *Repository) matchesFilter(part *repoModel.Part, filter *repoModel.PartsFilter) bool {
	if len(filter.PartUuids) > 0 && !contains(filter.PartUuids, part.UUID.String()) {
		return false
	}

	if len(filter.Names) > 0 && !contains(filter.Names, part.Name) {
		return false
	}

	if len(filter.Categories) > 0 && !contains(filter.Categories, part.Category) {
		return false
	}

	if len(filter.ManufacturerCountries) > 0 &&
		!contains(filter.ManufacturerCountries, part.Manufacturer.Country) {
		return false
	}

	if len(filter.Tags) > 0 && !hasAnyTag(part.Tags, filter.Tags) {
		return false
	}

	return true
}

func contains[T comparable](values []T, target T) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}

	return false
}

func hasAnyTag(partTags, filterTags []string) bool {
	for _, partTag := range partTags {
		for _, filterTag := range filterTags {
			if partTag == filterTag {
				return true
			}
		}
	}

	return false
}
