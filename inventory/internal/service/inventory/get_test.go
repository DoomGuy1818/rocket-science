package inventory

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/DoomGuy1818/rocket-science/inventory/internal/model"
)

func (s *ApiSuite) TestGetPart() {
	ctx := context.Background()

	dt := lo.ToPtr(time.Now())

	tests := []struct {
		name        string
		id          uuid.UUID
		err         error
		expectedErr error
		repoAns     model.Part
		expectedAns model.Part
	}{
		{
			name:        "Happy path with normal data",
			id:          uuid.New(),
			err:         nil,
			expectedErr: nil,
			repoAns: model.Part{
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
			expectedAns: model.Part{
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
		{
			name:        "Part not found case",
			id:          uuid.New(),
			err:         model.ErrNotFound,
			expectedErr: model.ErrNotFound,
			repoAns:     model.Part{},
			expectedAns: model.Part{},
		},
	}

	for _, tt := range tests {
		s.Run(
			tt.name, func() {
				s.repo.On("Get", mock.Anything, tt.id).
					Return(tt.repoAns, tt.err).
					Once()

				part, err := s.service.GetByID(ctx, tt.id)

				require.Equal(s.T(), tt.expectedErr, err)
				require.Equal(s.T(), tt.expectedAns, part)
			},
		)
	}
}
