package usecase

type UserYogaUsecase interface{}

type userYogaUseCase struct{}

func NewUserYogaUsecase() UserYogaUsecase {
	return &userYogaUseCase{}
}
