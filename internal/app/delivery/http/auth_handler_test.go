package http

import (
	"context"
	"fmt"
	"log"
	"net/http/httptest"
	"testing"

	"onthemat/internal/app/config"
	"onthemat/internal/app/infrastructor"
	"onthemat/internal/app/repository"
	"onthemat/internal/app/service"
	"onthemat/internal/app/service/token"
	"onthemat/internal/app/usecase"
	"onthemat/pkg/auth/jwt"
	"onthemat/pkg/kakao"

	"github.com/chromedp/chromedp"
	"github.com/gofiber/fiber/v2"
)

func Init() *fiber.App {
	configPath := "../../../../configs"

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
	authUseCase := usecase.NewAuthUseCase(tokenModule, userRepo, authSvc, c)
	userUseCase := usecase.NewUserUseCase(userRepo)

	// app
	app := fiber.New()
	router := app.Group("/api/v1")
	NewAuthHandler(authUseCase, userUseCase, router)

	return app
}

func TestAuth(t *testing.T) {
	app := Init()
	req := httptest.NewRequest("GET", "/api/v1/auth/kakao", nil)
	resp, _ := app.Test(req)
	fmt.Println(resp)
}

func click(ctx context.Context) {
	var example string
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://pkg.go.dev/time`),
		// wait for footer element is visible (ie, page is loaded)
		chromedp.WaitVisible(`body > footer`),
		// find and click "Example" link
		chromedp.Click(`#id_email_2`, chromedp.NodeVisible),
		// retrieve the text of the textarea
		chromedp.Value(`#example-After textarea`, &example),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Go's time.After example:\n%s", example)
}
