package fiberx

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func BodyParser(c *fiber.Ctx, out interface{}) error {
	err := c.BodyParser(out)

	e, ok := err.(*json.UnmarshalTypeError)
	if !ok {
		return err
	}

	message := fmt.Sprintf("%s field %s to %s", e.Field, e.Value, e.Type.Name())
	return errors.New(message)
}
