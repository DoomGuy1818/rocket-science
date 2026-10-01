package converter

import (
	"github.com/DoomGuy1818/rocket-science/inventory/internal/model"
	repoModel "github.com/DoomGuy1818/rocket-science/inventory/internal/repository/model"
)

func PartToRepo(part model.Part) repoModel.Part {
	return repoModel.Part{}
}

func PartToModel(part repoModel.Part) model.Part {
	return model.Part{}
}
