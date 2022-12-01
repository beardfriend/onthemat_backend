package middlewares

import (
	"onthemat/internal/app/repository"
	"onthemat/internal/app/service"
	"onthemat/internal/app/service/token"

	"github.com/gofiber/fiber/v2"
)

type MiddleWare interface {
	Auth(c *fiber.Ctx) error
	OnlyAcademy(c *fiber.Ctx) error
	OnlyTeacher(c *fiber.Ctx) error
	OnlySuperAdmin(c *fiber.Ctx) error
}

type middleWare struct {
	authSvc     service.AuthService
	tokensvc    token.TokenService
	teacherRepo repository.TeacherRepository
	academyRepo repository.AcademyRepository
}

func NewMiddelwWare(
	authSvc service.AuthService,
	tokensvc token.TokenService,
	teacherRepo repository.TeacherRepository,
	academyRepo repository.AcademyRepository,
) MiddleWare {
	return &middleWare{
		authSvc:     authSvc,
		tokensvc:    tokensvc,
		teacherRepo: teacherRepo,
		academyRepo: academyRepo,
	}
}
