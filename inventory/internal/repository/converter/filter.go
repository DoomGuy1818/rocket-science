package converter

import (
	"github.com/DoomGuy1818/rocket-science/inventory/internal/model"
	repoFilter "github.com/DoomGuy1818/rocket-science/inventory/internal/repository/model"
)

func ModelFilterToRepo(filter *model.PartsFilter) *repoFilter.PartsFilter {
	return &repoFilter.PartsFilter{
		PartUuids:             filter.PartUuids,
		Names:                 filter.Names,
		Categories:            modelCategoriesToRepo(filter.Categories),
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}

func modelCategoriesToRepo(categories []model.Category) []repoFilter.Category {
	repoCategories := make([]repoFilter.Category, len(categories))
	for _, category := range categories {
		repoCategories = append(repoCategories, repoFilter.Category(category))
	}

	return repoCategories
}
