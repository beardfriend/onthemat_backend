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
	"onthemat/pkg/auth/store/redis"
	"onthemat/pkg/email"
	"onthemat/pkg/google"
	"onthemat/pkg/kakao"
	"onthemat/pkg/naver"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
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
	jwt := jwt.NewJwt().WithSignKey(c.JWT.SignKey).Init()
	tokenModule := token.NewToken(jwt)
	k := kakao.NewKakao(c)
	g := google.NewGoogle(c)
	n := naver.NewNaver(c)
	emailM := email.NewEmail(c)

	// db
	db := infrastructor.NewPostgresDB()
	redisCli := infrastructor.NewRedis(c)

	// repo
	userRepo := repository.NewUserRepository(db)

	// service
	authSvc := service.NewAuthService(k, g, n, emailM)
	authStore := redis.NewStore(redisCli)

	// usecase
	authUseCase := usecase.NewAuthUseCase(tokenModule, userRepo, authSvc, authStore, c)
	userUsecase := usecase.NewUserUseCase(userRepo)

	// middleware
	middleWare := middlewares.NewMiddelwWare(authSvc, tokenModule)

	defer func() {
		infrastructor.ClosePostgres(db)
	}()
	// app
	app := fiber.New()
	app.Use(recover.New())

	// handler
	router := app.Group("/api/v1")
	http.NewAuthHandler(authUseCase, router)
	http.NewUserHandler(middleWare, userUsecase, router)

	app.Listen(":3000")
}
