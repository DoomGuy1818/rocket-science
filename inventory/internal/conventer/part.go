package conventer

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/DoomGuy1818/rocket-science/inventory/internal/model"
	inventoryV1 "github.com/DoomGuy1818/rocket-science/shared/pkg/proto/inventory/v1"
)

func InventoryPartToModel(part *inventoryV1.Part) (model.Part, error) {
	var createdAt *time.Time
	if part.CreatedAt != nil {
		createdAt = lo.ToPtr(part.CreatedAt.AsTime())
	}

	var updatedAt *time.Time
	if part.UpdatedAt != nil {
		updatedAt = lo.ToPtr(part.UpdatedAt.AsTime())
	}

	id, err := uuid.Parse(part.Uuid)
	if err != nil {
		return model.Part{}, err
	}

	return model.Part{
		UUID:          id,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: int(part.StockQuantity),
		Category:      categoryToModel(part.Category),
		Dimensions:    dimensionsToModel(part.Dimensions),
		Manufacturer:  manufacturerToModel(part.Manufacturer),
		Metadata:      metadataToModel(part.Metadata),
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}, nil
}

func ModelPartToProto(part model.Part) *inventoryV1.Part {
	var createdAt *timestamppb.Timestamp
	if part.CreatedAt != nil {
		createdAt = timestamppb.New(*part.CreatedAt)
	}

	var updatedAt *timestamppb.Timestamp
	if part.UpdatedAt != nil {
		updatedAt = timestamppb.New(*part.UpdatedAt)
	}

	return &inventoryV1.Part{
		Uuid:          part.UUID.String(),
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: int64(part.StockQuantity),
		Category:      categoryToProto(part.Category),
		Dimensions:    dimensionsToProto(part.Dimensions),
		Manufacturer:  manufacturerToProto(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      metadataToProto(part.Metadata),
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func ModelPartsToProto(parts []model.Part) []*inventoryV1.Part {
	prts := make([]*inventoryV1.Part, 0, len(parts))
	for _, part := range parts {
		prts = append(prts, ModelPartToProto(part))
	}

	return prts
}

func PartFilterToModel(filters *inventoryV1.PartFilters) (*model.PartsFilter, error) {
	uuids := make([]uuid.UUID, 0, len(filters.GetUuids()))
	for _, ID := range filters.Uuids {
		UUID, err := uuid.Parse(ID)
		if err != nil {
			return nil, fmt.Errorf("invalid uuid %q: %w", ID, err)
		}
		uuids = append(uuids, UUID)
	}

	var Categories []model.Category
	for _, category := range filters.Categories {
		Categories = append(Categories, categoryToModel(category))
	}

	return &model.PartsFilter{
		PartUuids:             uuids,
		Names:                 filters.Names,
		Categories:            Categories,
		ManufacturerCountries: filters.ManufacturerCountries,
		Tags:                  filters.Tags,
	}, nil
}

func categoryToModel(c inventoryV1.Category) model.Category {
	switch c {
	case inventoryV1.Category_ENGINE:
		return model.CategoryEngine
	case inventoryV1.Category_FUEL:
		return model.CategoryFuel
	case inventoryV1.Category_PORTHOLE:
		return model.CategoryPorthole
	case inventoryV1.Category_WING:
		return model.CategoryWing
	default:
		return model.CategoryUnknown
	}
}

func categoryToProto(c model.Category) inventoryV1.Category {
	switch c {
	case model.CategoryEngine:
		return inventoryV1.Category_ENGINE
	case model.CategoryFuel:
		return inventoryV1.Category_FUEL
	case model.CategoryPorthole:
		return inventoryV1.Category_PORTHOLE
	case model.CategoryWing:
		return inventoryV1.Category_WING
	default:
		return inventoryV1.Category_UNKNOWN
	}
}

func dimensionsToModel(d *inventoryV1.Dimensions) model.Dimensions {
	return model.Dimensions{
		Length: float32(d.Length),
		Width:  float32(d.Width),
		Height: float32(d.Height),
		Weight: float32(d.Weight),
	}
}

func manufacturerToModel(m *inventoryV1.Manufacturer) model.Manufacturer {
	return model.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

func metadataToModel(meta map[string]*inventoryV1.Value) model.Metadata {
	if meta == nil {
		return nil
	}

	result := make(model.Metadata, len(meta))

	for k, v := range meta {
		switch val := v.GetValue().(type) {
		case *inventoryV1.Value_StringValue:
			result[k] = model.StringValue(val.StringValue)
		case *inventoryV1.Value_Int64Value:
			result[k] = model.Int64Value(val.Int64Value)
		case *inventoryV1.Value_DoubleValue:
			result[k] = model.DoubleValue(val.DoubleValue)
		case *inventoryV1.Value_BoolValue:
			result[k] = model.BoolValue(val.BoolValue)
		default:
		}
	}
	return result
}

func dimensionsToProto(d model.Dimensions) *inventoryV1.Dimensions {
	return &inventoryV1.Dimensions{
		Length: float64(d.Length),
		Width:  float64(d.Width),
		Height: float64(d.Height),
		Weight: float64(d.Weight),
	}
}

func manufacturerToProto(m model.Manufacturer) *inventoryV1.Manufacturer {
	return &inventoryV1.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
	}
}

func metadataToProto(meta model.Metadata) map[string]*inventoryV1.Value {
	if meta == nil {
		return nil
	}

	result := make(map[string]*inventoryV1.Value, len(meta))

	for k, v := range meta {
		switch val := v.(type) {
		case model.StringValue:
			result[k] = &inventoryV1.Value{Value: &inventoryV1.Value_StringValue{StringValue: string(val)}}
		case model.Int64Value:
			result[k] = &inventoryV1.Value{Value: &inventoryV1.Value_Int64Value{Int64Value: int64(val)}}
		case model.DoubleValue:
			result[k] = &inventoryV1.Value{Value: &inventoryV1.Value_DoubleValue{DoubleValue: float64(val)}}
		case model.BoolValue:
			result[k] = &inventoryV1.Value{Value: &inventoryV1.Value_BoolValue{BoolValue: bool(val)}}
		default:
			// nil-значение в map: пропускаем
		}
	}
	return result
}
