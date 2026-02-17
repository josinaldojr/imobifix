package handlers

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/josinaldojr/imobifix-api/internal/errors"
	"github.com/josinaldojr/imobifix-api/internal/integrations/viacep"
)

var cepRe = regexp.MustCompile(`^\d{8}$`)

func normalizeCEP(s string) (string, bool) {
	d := strings.NewReplacer("-", "", ".", "", " ", "").Replace(s)
	return d, cepRe.MatchString(d)
}

func Address(client *viacep.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		raw := c.Params("cep")
		cep, ok := normalizeCEP(raw)
		if !ok {
			return errors.New(http.StatusBadRequest, "CEP_INVALID", "CEP inválido. Use 8 dígitos.", fiber.Map{"cep": raw})
		}

		addr, err := client.Lookup(context.Background(), cep)
		if err == viacep.ErrNotFound {
			return errors.New(http.StatusNotFound, "CEP_NOT_FOUND", "CEP não encontrado.", fiber.Map{"cep": raw})
		}
		if err == viacep.ErrUnavailable || err == viacep.ErrInvalidAnswer {
			return errors.New(http.StatusServiceUnavailable, "VIA_CEP_UNAVAILABLE", "Não foi possível consultar o CEP agora. Preencha o endereço manualmente.", nil)
		}
		if err != nil {
			return err
		}
		return c.JSON(addr)
	}
}
