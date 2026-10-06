package converter

import (
	"github.com/DoomGuy1818/rocket-science/inventory/internal/model"
	repoModel "github.com/DoomGuy1818/rocket-science/inventory/internal/repository/model"
)

func PartToRepo(part model.Part) repoModel.Part {
	return repoModel.Part{
		UUID:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      modelCategoryToRepo(part.Category),
		Dimensions:    modelDimensionsToRepo(part.Dimensions),
		Manufacturer:  modelManufacturerToRepo(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      modelMetadataToRepo(part.Metadata),
	}
}

func PartToModel(part repoModel.Part) model.Part {
	return model.Part{
		UUID:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      repoCategoryToModel(part.Category),
		Dimensions:    repoDimensionsToModel(part.Dimensions),
		Manufacturer:  repoManufacturerToModel(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      repoMetadataToModel(part.Metadata),
	}
}

func modelCategoryToRepo(category model.Category) repoModel.Category {
	return repoModel.Category(category)
}

func modelDimensionsToRepo(dimensions model.Dimensions) repoModel.Dimensions {
	return repoModel.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func modelManufacturerToRepo(manufacturer model.Manufacturer) repoModel.Manufacturer {
	return repoModel.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
	}
}

func modelMetadataToRepo(metadata model.Metadata) repoModel.Metadata {
	result := make(repoModel.Metadata, len(metadata))

	for k, v := range metadata {
		switch val := v.(type) {
		case model.StringValue:
			result[k] = repoModel.StringValue(val)
		case model.Int64Value:
			result[k] = repoModel.Int64Value(val)
		case model.DoubleValue:
			result[k] = repoModel.DoubleValue(val)
		case model.BoolValue:
			result[k] = repoModel.BoolValue(val)
		default:

		}
	}
	return result
}

func repoCategoryToModel(category repoModel.Category) model.Category {
	return model.Category(category)
}

func repoDimensionsToModel(dimensions repoModel.Dimensions) model.Dimensions {
	return model.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func repoManufacturerToModel(manufacturer repoModel.Manufacturer) model.Manufacturer {
	return model.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func repoMetadataToModel(metadata repoModel.Metadata) model.Metadata {
	if metadata == nil {
		return nil
	}

	result := make(model.Metadata, len(metadata))

	for k, v := range metadata {
		switch val := v.(type) {
		case repoModel.StringValue:
			result[k] = model.StringValue(val)
		case repoModel.Int64Value:
			result[k] = model.Int64Value(val)
		case repoModel.DoubleValue:
			result[k] = model.DoubleValue(val)
		case repoModel.BoolValue:
			result[k] = model.BoolValue(val)
		default:

		}
	}
	return result
}
