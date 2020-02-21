package usecase

import (
	"github.com/samaita/go-default/src"
)

type Usecase struct {
	repos src.Repository
}

func NewUsecase(ur src.Repository) src.Usecase {
	return &Usecase{
		repos: ur,
	}
}

func init() {
	// gob.Register(models.User{})
}
