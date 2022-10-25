package main

import (
	"flag"

	"onthemat/internal/app/config"
	"onthemat/internal/app/delivery/http"
	"onthemat/internal/app/infrastructor"
	"onthemat/internal/app/repository"
	"onthemat/internal/app/service"
	"onthemat/internal/app/service/token"
	"onthemat/internal/app/usecase"
	"onthemat/pkg/auth/jwt"
	"onthemat/pkg/kakao"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// debugmode check
	configPath := "./configs"
	isDebug := flag.Bool("mode", false, "debug")
	flag.Parse()
	if *isDebug {
		configPath = "../../configs"
	}
	// config
	c := config.NewConfig()
	if err := c.Load(configPath); err != nil {
		panic(err)
	}
	// pkg
	jwt := jwt.NewJwt().Init()
	tokenModule := token.NewToken(jwt)
	k := kakao.NewKakao(c)

	// db
	db := infrastructor.NewPostgresDB()

	// repo
	userRepo := repository.NewUserRepository(db)

	// service
	authSvc := service.NewAuthService(k)

	// usecase
	authUseCase := usecase.NewAuthUseCase(tokenModule, userRepo, authSvc)

	// app
	app := fiber.New()

	// handler
	router := app.Group("/api/v1")
	http.NewAuthHandler(authUseCase, router)

	app.Listen(":3000")
}
