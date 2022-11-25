package main

import (
	"context"
	"fmt"

	"onthemat/internal/app/config"
	"onthemat/internal/app/infrastructure"
	"onthemat/internal/app/repository"
	"onthemat/internal/app/service"
	"onthemat/internal/app/usecase"
	"onthemat/internal/app/utils"
)

func main() {
	seed := utils.NewSeeding()
	seed.Users()
	c := config.NewConfig()
	if err := c.Load("./configs"); err != nil {
		panic(err)
	}
	db := infrastructure.NewPostgresDB(c)
	areaRepo := repository.NewAreaRepository(db)
	areaService := service.NewAreaService()
	areaUsecase := usecase.NewAreaUsecase(areaRepo, areaService)
	err := areaUsecase.CreateSiDo(context.Background(), "/Users/sehun/Downloads/행정.xlsx")
	fmt.Println(err)
	seed.Academies()
	seed.YogaGroup()
	seed.Yoga()
}
