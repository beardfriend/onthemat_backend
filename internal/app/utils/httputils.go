package utils

import (
	"onthemat/internal/app/common"

	"github.com/gofiber/fiber/v2"
)

func NewError(c *fiber.Ctx, err error) error {
	code, json := common.ParseHttpError(err)
	return c.Status(code).JSON(json)
}
