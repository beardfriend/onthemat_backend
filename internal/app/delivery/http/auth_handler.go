package http

import (
	"context"

	"onthemat/internal/app/dto"
	"onthemat/internal/app/repository"
	"onthemat/pkg/ent"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler interface {
	SignUpTest(c *fiber.Ctx) error
}

type authHandler struct {
	AcademyRepo repository.AcademyRepository
}

func NewAuthHandler(academyRepo repository.AcademyRepository) AuthHandler {
	return &authHandler{
		AcademyRepo: academyRepo,
	}
}

func (h *authHandler) SignUpTest(c *fiber.Ctx) error {
	context := context.Background()
	defer context.Done()

	body := new(dto.AcademyNormalSignUpBody)
	if err := c.BodyParser(body); err != nil {
		return err
	}
	res := &ent.Acadmey{
		Name:         body.NickName,
		BusinessCode: &body.BusinessCode,
		FullAddress:  body.Address,

		Edges: ent.AcadmeyEdges{
			User: &ent.User{
				Email:    body.Email,
				Password: body.Password,
				Nickname: body.NickName,
			},
		},
	}
	if err := h.AcademyRepo.Create(context, res); err != nil {
		return err
	}

	return c.SendStatus(200)
}
