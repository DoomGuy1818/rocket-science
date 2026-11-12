package inventory

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/DoomGuy1818/rocket-science/inventory/internal/model"
)

func (i *repository) List(ctx context.Context, filters *model.PartsFilter) ([]model.Part, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	cursor, err := i.parts.Find(
		ctx, buildFilter(filters),
	)
	if err != nil {
		return nil, err
	}

	defer func() {
		cerr := cursor.Close(ctx)
		if cerr != nil {
			log.Printf("failed to close cursor: %v\n", cerr)
		}
	}()

	var parts []model.Part

	err = cursor.All(ctx, &parts)
	if err != nil {
		return nil, err
	}

	for _, part := range parts {
		log.Printf(
			"Part: %s, Metadata: %#v",
			part.Name,
			part.Metadata,
		)
	}

	return parts, nil
}

func buildFilter(filters *model.PartsFilter) bson.M {
	filter := bson.M{}

	if len(filters.PartUuids) > 0 {
		filter["uuid"] = bson.M{
			"$in": filters.PartUuids,
		}
	}

	if len(filters.Names) > 0 {
		filter["name"] = bson.M{
			"$in": filters.Names,
		}
	}

	if len(filters.Categories) > 0 {
		filter["category"] = bson.M{
			"$in": filters.Categories,
		}
	}

	if len(filters.ManufacturerCountries) > 0 {
		filter["manufacturerCountries"] = bson.M{
			"$in": filters.ManufacturerCountries,
		}
	}

	if len(filters.Tags) > 0 {
		filter["tags"] = bson.M{
			"$in": filters.Tags,
		}
	}

	return filter
}
