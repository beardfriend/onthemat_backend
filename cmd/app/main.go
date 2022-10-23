package main

import (
	"onthemat/internal/app/delivery/http"
	"onthemat/internal/app/infrastructor"
	"onthemat/internal/app/repository"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	db := infrastructor.NewPostgresDB()
	repo := repository.NewAcademyRepository(db)
	handler := http.NewAuthHandler(repo)

	app.Post("/", handler.SignUpTest)
	app.Listen(":3000")
}
