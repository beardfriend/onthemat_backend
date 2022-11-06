package main

import (
	"flag"

	"onthemat/internal/app/config"
	"onthemat/internal/app/delivery/http"
	"onthemat/internal/app/delivery/middlewares"
	"onthemat/internal/app/infrastructor"
	"onthemat/internal/app/repository"
	"onthemat/internal/app/service"
	"onthemat/internal/app/service/token"
	"onthemat/internal/app/usecase"
	"onthemat/pkg/auth/jwt"
	"onthemat/pkg/google"
	"onthemat/pkg/kakao"
	"onthemat/pkg/naver"

	swagger "github.com/arsmn/fiber-swagger/v2"

	_ "onthemat/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// @title Go Boilerplate
// @version 1.0.4
// @termsOfService http://swagger.io/terms/
// @securityDefinitions.apikey JwtToken
// @in header
// @name Authorization
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
	jwt := jwt.NewJwt().WithSignKey(c.JWT.SignKey).Init()
	tokenModule := token.NewToken(jwt)
	k := kakao.NewKakao(c)
	g := google.NewGoogle(c)
	n := naver.NewNaver(c)

	// db
	db := infrastructor.NewPostgresDB()

	// repo
	userRepo := repository.NewUserRepository(db)

	// service
	authSvc := service.NewAuthService(k, g, n)

	// usecase
	authUseCase := usecase.NewAuthUseCase(tokenModule, userRepo, authSvc, c)
	userUsecase := usecase.NewUserUseCase(userRepo)

	// middleware

	middleWare := middlewares.NewMiddelwWare(authSvc, tokenModule)

	defer func() {
		infrastructor.ClosePostgres(db)
	}()
	// app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(500).JSON(fiber.Map{
				"message": "일시적인 에러가 발생했습니다",
			})
		},
	})
	app.Use(recover.New())
	app.Get("/swagger/*", swagger.HandlerDefault)

	// handler
	router := app.Group("/api/v1")
	http.NewAuthHandler(authUseCase, router)
	http.NewUserHandler(middleWare, userUsecase, router)

	app.Listen(":3000")
}
