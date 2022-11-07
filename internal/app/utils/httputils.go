package utils

import (
	"onthemat/internal/app/common"

	"github.com/gofiber/fiber/v2"
)

func NewError(c *fiber.Ctx, status int, err common.HttpErr) error {
	return c.Status(status).JSON(err)
}
