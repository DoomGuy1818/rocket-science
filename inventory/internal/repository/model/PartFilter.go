package model

type PartsFilter struct {
	PartUuids             []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}
