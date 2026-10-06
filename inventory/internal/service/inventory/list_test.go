package inventory

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/DoomGuy1818/rocket-science/inventory/internal/model"
)

func (s *ApiSuite) TestListParts() {
	ctx := context.Background()

	dt := lo.ToPtr(time.Now())

	tests := []struct {
		name     string
		filters  model.PartsFilter
		repoAns  []model.Part
		expParts []model.Part
	}{
		{
			name: "successful get some parts",
			filters: model.PartsFilter{
				PartUuids: []uuid.UUID{uuid.MustParse("3f2b8c1e-7a4d-4e9b-9c3a-5d6e8f1a2b47")},
			},
			repoAns: []model.Part{
				{
					UUID:          uuid.MustParse("3f2b8c1e-7a4d-4e9b-9c3a-5d6e8f1a2b47"),
					Name:          "Ионный двигатель X-200",
					Description:   "Маломощный ионный двигатель для орбитальных манёвров",
					Price:         1250000.00,
					StockQuantity: 5,
					Category:      model.CategoryEngine,
					Dimensions:    model.Dimensions{Length: 220.5, Width: 90, Height: 90, Weight: 340},
					Manufacturer:  model.Manufacturer{Name: "AeroDyne", Country: "USA"},
					Tags:          []string{"ion", "orbital", "low-thrust"},
					Metadata: model.Metadata{
						"cert_number": model.StringValue("ENG-2024-0042"),
					},
					CreatedAt: dt,
					UpdatedAt: dt,
				},
			},
			expParts: []model.Part{
				{
					UUID:          uuid.MustParse("3f2b8c1e-7a4d-4e9b-9c3a-5d6e8f1a2b47"),
					Name:          "Ионный двигатель X-200",
					Description:   "Маломощный ионный двигатель для орбитальных манёвров",
					Price:         1250000.00,
					StockQuantity: 5,
					Category:      model.CategoryEngine,
					Dimensions:    model.Dimensions{Length: 220.5, Width: 90, Height: 90, Weight: 340},
					Manufacturer:  model.Manufacturer{Name: "AeroDyne", Country: "USA"},
					Tags:          []string{"ion", "orbital", "low-thrust"},
					Metadata: model.Metadata{
						"cert_number": model.StringValue("ENG-2024-0042"),
					},
					CreatedAt: dt,
					UpdatedAt: dt,
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(
			test.name, func() {
				s.repo.On("List", ctx, &test.filters).Return(test.repoAns).Once()

				s.service.ListByFilters(ctx, &test.filters)
			},
		)
	}
}
