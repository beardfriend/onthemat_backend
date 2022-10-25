package main

import (
	"onthemat/internal/app/config"
	"onthemat/internal/app/delivery/http"
	"onthemat/internal/app/infrastructor"
	"onthemat/internal/app/repository"
	"onthemat/internal/app/service/token"
	"onthemat/internal/app/usecase"
	"onthemat/pkg/auth/jwt"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// config
	c := config.NewConfig()
	// if run  "./configs"
	if err := c.Load("../../configs"); err != nil {
		panic(err)
	}
	// service
	jwt := jwt.NewJwt().Init()
	tokenModule := token.NewToken(jwt)

	// db
	db := infrastructor.NewPostgresDB()

	// repo
	userRepo := repository.NewUserRepository(db)

	// usecase
	authUseCase := usecase.NewAuthUseCase(c, tokenModule, userRepo)

	// handler
	router := app.Group("/api/v1")
	http.NewAuthHandler(authUseCase, router)

	app.Listen(":3000")
}
