package main

import (
	"context"

	"onthemat/internal/app/config"
	"onthemat/internal/app/infrastructure"
	"onthemat/internal/app/repository"
	"onthemat/internal/app/service"
	"onthemat/internal/app/usecase"
)

func main() {
	c := config.NewConfig()
	if err := c.Load("./configs"); err != nil {
		panic(err)
	}
	db := infrastructure.NewPostgresDB(c)
	areaRepo := repository.NewAreaRepository(db)
	areaService := service.NewAreaService()
	areaUsecase := usecase.NewAreaUsecase(areaRepo, areaService)
	areaUsecase.CreateSiDo(context.Background(), "/Users/sehun/Downloads/국토교통부.xlsx")
}
