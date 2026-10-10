package model

import "github.com/google/uuid"

type PartsFilter struct {
	PartUuids             []uuid.UUID
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}
