package main

import (
	"flag"

	"onthemat/internal/app/config"
	"onthemat/internal/app/delivery/http"
	"onthemat/internal/app/delivery/middlewares"
	"onthemat/internal/app/infrastructure"
	"onthemat/internal/app/repository"
	"onthemat/internal/app/service"
	"onthemat/internal/app/service/token"
	"onthemat/internal/app/usecase"
	"onthemat/pkg/auth/jwt"
	"onthemat/pkg/auth/store/redis"
	"onthemat/pkg/aws"
	"onthemat/pkg/email"
	"onthemat/pkg/google"
	"onthemat/pkg/kakao"
	"onthemat/pkg/naver"
	"onthemat/pkg/openapi"
	"onthemat/pkg/validatorx"

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

	// ------- s3 ----------
	s3 := aws.NewS3(c)
	if err := s3.SetConfig(); err != nil {
		panic(err)
	}

	businessManM := openapi.NewBusinessMan(c)

	// db
	db := infrastructure.NewPostgresDB(c)
	redisCli := infrastructure.NewRedis(c)

	// utils
	validator := validatorx.NewValidatorx().
		AddPasswordAtLeastOneCharNumValidation("PassWordAtLeastOneCharOneNum").
		AddPhoneNumValidation("phoneNumNoDash").
		AddUrlValidation("urlStartHttpHttps").Init()

	// repo
	userRepo := repository.NewUserRepository(db)
	imageRepo := repository.NewImageRepository(db)
	academyRepo := repository.NewAcademyRepository(db)

	// service
	authSvc := service.NewAuthService(k, g, n, emailM)
	authStore := redis.NewStore(redisCli)
	academySvc := service.NewAcademyService(businessManM)

	// usecase
	authUseCase := usecase.NewAuthUseCase(tokenModule, userRepo, authSvc, authStore, c)
	userUsecase := usecase.NewUserUseCase(userRepo)
	academyUsecase := usecase.NewAcademyUsecase(academyRepo, academySvc, userRepo)
	uploadUsecase := usecase.NewUploadUsecase(imageRepo, s3)

	// middleware
	middleWare := middlewares.NewMiddelwWare(authSvc, tokenModule)

	defer func() {
		infrastructure.ClosePostgres(db)
	}()
	// app
	app := fiber.New()
	app.Use(recover.New())
	// app.Use(limiter.New(limiter.ConfigDefault))

	// handler
	router := app.Group("/api/v1")
	http.NewAuthHandler(authUseCase, validator, router)
	http.NewUploadHandler(middleWare, uploadUsecase, validator, router)
	http.NewUserHandler(middleWare, userUsecase, router)
	http.NewAcademyHandler(middleWare, academyUsecase, validator, router)

	app.Listen(":3000")
}
