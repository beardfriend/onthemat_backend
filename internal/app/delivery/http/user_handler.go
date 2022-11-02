package http

import (
	"onthemat/internal/app/delivery/middlewares"
	"onthemat/internal/app/transport"
	"onthemat/internal/app/usecase"

	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
	UserUseCase usecase.UserUseCase
	Router      fiber.Router
	Middleware  *middlewares.MiddleWare
}

func NewUserHandler(
	middleware *middlewares.MiddleWare,
	userUseCase usecase.UserUseCase,
	router fiber.Router,
) {
	handler := &userHandler{
		UserUseCase: userUseCase,
		Middleware:  middleware,
	}
	g := router.Group("/user")
	g.Get("/me", middleware.Auth, handler.GetMe)
}

func (h *userHandler) GetMe(c *fiber.Ctx) error {
	ctx := c.Context()

	u, e := h.UserUseCase.GetMe(ctx, c.Context().UserValue("user_id").(int))

	if e != nil {
		panic(e)
	}

	resp := transport.NewUserMeResponse(u)
	return c.JSON(resp)
}
