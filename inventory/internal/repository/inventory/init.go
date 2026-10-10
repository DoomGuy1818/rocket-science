package inventory

import (
	"time"

	"github.com/google/uuid"

	"github.com/DoomGuy1818/rocket-science/inventory/internal/repository/model"
)

func InitRepository() map[uuid.UUID]*model.Part {
	now := time.Now().UTC()

	ts := func(year int, month time.Month, day, hour, min int) *time.Time {
		t := time.Date(year, month, day, hour, min, 0, 0, time.UTC)
		return &t
	}

	list := []*model.Part{
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
			CreatedAt: ts(2024, time.January, 15, 10, 30),
			UpdatedAt: ts(2024, time.June, 1, 8, 0),
		},
		{
			UUID:          uuid.MustParse("a81d5e6c-2f94-4b7a-8e13-c07b9d4f6a20"),
			Name:          "Химический двигатель R-90",
			Description:   "Жидкостный ракетный двигатель для первой ступени",
			Price:         3400000.50,
			StockQuantity: 2,
			Category:      model.CategoryEngine,
			Dimensions:    model.Dimensions{Length: 310, Width: 140, Height: 140, Weight: 1200},
			Manufacturer:  model.Manufacturer{Name: "Energomash", Country: "Russia"},
			Tags:          []string{"liquid", "first-stage", "heavy"},
			Metadata: model.Metadata{
				"thrust_kn": model.Int64Value(1950),
			},
			CreatedAt: ts(2023, time.September, 12, 9, 0),
			UpdatedAt: nil,
		},
		{
			UUID:          uuid.MustParse("9c4e7b2a-d1f8-4a35-b6e9-12a3c8d5f704"),
			Name:          "Топливный бак T-50",
			Description:   "Титановый бак на 5000 литров для жидкого топлива",
			Price:         480000.50,
			StockQuantity: 12,
			Category:      model.CategoryFuel,
			Dimensions:    model.Dimensions{Length: 300, Width: 150, Height: 150, Weight: 820},
			Manufacturer:  model.Manufacturer{Name: "CryoTank GmbH", Country: "Germany"},
			Tags:          []string{"titanium", "cryogenic"},
			Metadata: model.Metadata{
				"pressure_bar": model.DoubleValue(3.5),
			},
			CreatedAt: ts(2024, time.February, 3, 12, 0),
			UpdatedAt: ts(2024, time.April, 18, 15, 10),
		},
		{
			UUID:          uuid.MustParse("e5b0f3d7-6c28-4d91-a74e-8b2f1c9d3e65"),
			Name:          "Иллюминатор P-12",
			Description:   "Трёхслойный иллюминатор с защитой от микрометеоритов",
			Price:         75000.00,
			StockQuantity: 40,
			Category:      model.CategoryPorthole,
			Dimensions:    model.Dimensions{Length: 45, Width: 45, Height: 12, Weight: 18.5},
			Manufacturer:  model.Manufacturer{Name: "Clearview Optics", Country: "Japan"},
			Tags:          []string{"glass", "radiation-shield", "triple-layer"},
			Metadata: model.Metadata{
				"uv_filter": model.BoolValue(true),
			},
			CreatedAt: ts(2023, time.November, 20, 9, 15),
			UpdatedAt: ts(2024, time.March, 10, 14, 45),
		},
		{
			UUID:          uuid.MustParse("1d7a9f4c-b3e6-4258-8f0d-6a5c2e7b9314"),
			Name:          "Крыло W-7 Delta",
			Description:   "Дельтовидное крыло из углепластика для атмосферного спуска",
			Price:         920000.00,
			StockQuantity: 0, // нет в наличии
			Category:      model.CategoryWing,
			Dimensions:    model.Dimensions{Length: 650, Width: 400, Height: 35, Weight: 510},
			Manufacturer:  model.Manufacturer{Name: "SkyForge", Country: "France"},
			Tags:          []string{"carbon", "delta", "reentry"},
			Metadata: model.Metadata{
				"max_temp_c": model.Int64Value(1650),
			},
			CreatedAt: ts(2023, time.August, 5, 7, 0),
			UpdatedAt: ts(2024, time.May, 22, 16, 20),
		},
		{
			UUID:          uuid.MustParse("c62f8a1b-94d7-4e03-b5a8-f3d19e7c4b82"),
			Name:          "Крыло W-3 Straight",
			Description:   "Прямое крыло для лёгких суборбитальных аппаратов",
			Price:         310000.00,
			StockQuantity: 7,
			Category:      model.CategoryWing,
			Dimensions:    model.Dimensions{Length: 420, Width: 180, Height: 22, Weight: 240},
			Manufacturer:  model.Manufacturer{Name: "SkyForge", Country: "France"},
			Tags:          []string{"aluminium", "suborbital"},
			Metadata: model.Metadata{
				"foldable": model.BoolValue(true),
			},
			CreatedAt: ts(2024, time.January, 9, 11, 0),
			UpdatedAt: nil,
		},
		{
			// Пограничный случай: пустые значения и nil-ы
			UUID:          uuid.MustParse("7b3e5d90-a2c1-4f68-9d47-e0b8a6c13f5a"),
			Name:          "Неопознанный модуль",
			Description:   "Деталь без категории, для проверки edge-кейсов",
			Price:         0,
			StockQuantity: 1,
			Category:      model.CategoryUnknown,
			Dimensions:    model.Dimensions{},
			Manufacturer:  model.Manufacturer{},
			Tags:          nil,
			Metadata:      nil,
			CreatedAt:     &now,
			UpdatedAt:     nil,
		},
	}

	parts := make(map[uuid.UUID]*model.Part, len(list))
	for _, p := range list {
		parts[p.UUID] = p
	}

	return parts
}
