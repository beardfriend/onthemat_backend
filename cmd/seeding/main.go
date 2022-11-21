package main

import (
	"onthemat/internal/app/utils"
)

func main() {
	seed := utils.NewSeeding()
	seed.Users()
	seed.Academies()
	// c := config.NewConfig()
	// if err := c.Load("./configs"); err != nil {
	// 	panic(err)
	// }
	// db := infrastructure.NewPostgresDB(c)
	// areaRepo := repository.NewAreaRepository(db)
	// areaService := service.NewAreaService()
	// areaUsecase := usecase.NewAreaUsecase(areaRepo, areaService)
	// areaUsecase.CreateSiDo(context.Background(), "/Users/sehun/Downloads/행정.xlsx")
}
