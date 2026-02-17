package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/josinaldojr/imobifx-api/internal/service"
)

func Address(svc *service.AddressService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		addr, err := svc.Lookup(c.UserContext(), c.Params("cep"))
		if err != nil {
			return err
		}
		return c.JSON(addr)
	}
}
