package usecase

import "context"

type HealthUseCase interface{}

type healthUseCase struct{}

func NewHealthUseCase() HealthUseCase {
	return &healthUseCase{}
}

func (healthUseCase) Get(c context.Context) {
}
