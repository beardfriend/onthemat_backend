package main

import (
	"flag"
	"fmt"
	"log"
	"reflect"
	"time"

	"onthemat/internal/app/config"
	"onthemat/internal/app/delivery/http"
	"onthemat/internal/app/delivery/middlewares"
	"onthemat/internal/app/infrastructure"
	"onthemat/internal/app/repository"
	"onthemat/internal/app/service"
	"onthemat/internal/app/service/token"
	"onthemat/internal/app/transport"
	"onthemat/internal/app/usecase"
	"onthemat/pkg/auth/jwt"
	"onthemat/pkg/auth/store/redis"
	"onthemat/pkg/aws"
	ela "onthemat/pkg/elastic"
	"onthemat/pkg/email"
	"onthemat/pkg/google"
	"onthemat/pkg/kakao"
	"onthemat/pkg/naver"
	"onthemat/pkg/openapi"
	"onthemat/pkg/validatorx"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
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
		log.Fatal(err)
	}

	businessManM := openapi.NewBusinessMan(c)

	// db
	db := infrastructure.NewPostgresDB(c)
	redisCli := infrastructure.NewRedis(c)
	elastic := infrastructure.NewElasticSearch(c)
	ela.InitYoga(elastic)
	// utils
	validator := validatorx.NewValidatorx().
		AddPasswordAtLeastOneCharNumValidation("PassWordAtLeastOneCharOneNum").
		AddPhoneNumValidation("phoneNumNoDash").
		AddCheckMustFieldIfIdFieldExistValidation("must").
		AddUrlValidation("urlStartHttpHttps").SetExtractTagName().Init()

	// repo
	userRepo := repository.NewUserRepository(db)
	imageRepo := repository.NewImageRepository(db)
	academyRepo := repository.NewAcademyRepository(db)
	areaRepo := repository.NewAreaRepository(db)
	yogaRepo := repository.NewYogaRepository(db, elastic)
	teacherRepo := repository.NewTeacherRepository(db)
	recruitmentRepo := repository.NewRecruitmentRepository(db)

	// service
	authSvc := service.NewAuthService(k, g, n, emailM)
	authStore := redis.NewStore(redisCli)
	academySvc := service.NewAcademyService(businessManM)

	// usecase
	authUseCase := usecase.NewAuthUseCase(tokenModule, userRepo, authSvc, authStore, c)
	userUsecase := usecase.NewUserUseCase(userRepo)
	academyUsecase := usecase.NewAcademyUsecase(academyRepo, academySvc, userRepo, yogaRepo, areaRepo)
	uploadUsecase := usecase.NewUploadUsecase(imageRepo, s3)
	yogaUsecase := usecase.NewYogaUsecase(yogaRepo, academyRepo, teacherRepo)
	teacherUsecase := usecase.NewTeacherUsecase(teacherRepo, userRepo)
	recruitmentUsecase := usecase.NewRecruitmentUsecase(recruitmentRepo)
	// middleware
	middleWare := middlewares.NewMiddelwWare(authSvc, tokenModule, teacherRepo, academyRepo)

	defer func() {
		infrastructure.ClosePostgres(db)
	}()
	// app

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	fiber.SetParserDecoder(fiber.ParserConfig{
		ParserType: []fiber.ParserType{{
			Customtype: transport.TimeString{},
			Converter: func(value string) reflect.Value {
				fmt.Println("timeConverter", value)
				if v, err := time.Parse("2006-01-02T15:04:05", value); err == nil {
					return reflect.ValueOf(v)
				}
				return reflect.Value{}
			},
		}},
	})

	app.Use(logger.New())
	app.Use(recover.New())
	// app.Use(limiter.New(limiter.ConfigDefault))

	// handler
	router := app.Group("/api/v1")
	http.NewAuthHandler(authUseCase, validator, router)
	http.NewUploadHandler(middleWare, uploadUsecase, validator, router)
	http.NewUserHandler(middleWare, userUsecase, router)
	http.NewAcademyHandler(middleWare, academyUsecase, validator, router)
	http.NewYogaHandler(yogaUsecase, middleWare, validator, router)
	http.NewTeacherHandler(middleWare, teacherUsecase, validator, router)
	http.NewRecruitmentHandler(middleWare, recruitmentUsecase, validator, router)
	app.Listen(":3000")
}
