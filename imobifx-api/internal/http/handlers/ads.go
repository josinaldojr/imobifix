package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/josinaldojr/imobifix-api/internal/config"
	"github.com/josinaldojr/imobifix-api/internal/domain"
	"github.com/josinaldojr/imobifix-api/internal/errors"
	"github.com/josinaldojr/imobifix-api/internal/repo"
	"github.com/josinaldojr/imobifix-api/internal/service"
)

func CreateAd(cfg config.Config, ads *service.AdsService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		typ := strings.ToUpper(strings.TrimSpace(c.FormValue("type")))
		priceStr := strings.TrimSpace(c.FormValue("price_brl"))

		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil || price < 0 {
			return errors.New(http.StatusBadRequest, "VALIDATION_ERROR", "Dados inválidos.", fiber.Map{"price_brl": "must be a number >= 0"})
		}

		cepRaw := c.FormValue("cep")
		cep, ok := normalizeCEP(cepRaw)
		
		if !ok {
			return errors.New(http.StatusBadRequest, "VALIDATION_ERROR", "Dados inválidos.", fiber.Map{"cep": "must have 8 digits"})
		}

		street := strings.TrimSpace(c.FormValue("street"))
		neigh := strings.TrimSpace(c.FormValue("neighborhood"))
		city := strings.TrimSpace(c.FormValue("city"))
		state := strings.ToUpper(strings.TrimSpace(c.FormValue("state")))

		details := fiber.Map{}
		if typ != "SALE" && typ != "RENT" {
			details["type"] = "must be SALE or RENT"
		}
		if street == "" {
			details["street"] = "required"
		}
		if neigh == "" {
			details["neighborhood"] = "required"
		}
		if city == "" {
			details["city"] = "required"
		}
		if len(state) != 2 {
			details["state"] = "must have 2 letters (UF)"
		}
		if len(details) > 0 {
			return errors.New(http.StatusBadRequest, "VALIDATION_ERROR", "Dados inválidos.", details)
		}

		var number *string
		if v := strings.TrimSpace(c.FormValue("number")); v != "" {
			number = &v
		}
		var complement *string
		if v := strings.TrimSpace(c.FormValue("complement")); v != "" {
			complement = &v
		}

		var imageName *string
		file, ferr := c.FormFile("image")
		if ferr == nil && file != nil {
			if file.Size > cfg.MaxImageBytes {
				return errors.New(http.StatusBadRequest, "IMAGE_TOO_LARGE", "Imagem excede o tamanho máximo.", fiber.Map{"max_bytes": cfg.MaxImageBytes})
			}

			ct := strings.ToLower(file.Header.Get("Content-Type"))
			if ct != "image/jpeg" && ct != "image/png" && ct != "image/webp" {
				return errors.New(http.StatusBadRequest, "UNSUPPORTED_IMAGE_TYPE", "Tipo de imagem não suportado (jpeg/png/webp).", fiber.Map{"content_type": ct})
			}

			src, err := file.Open()
			if err != nil {
				return err
			}
			defer src.Close()

			name := ads.SaveImage(file.Filename)
			dstPath := ads.ImagePath(name)

			if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
				return err
			}

			dst, err := os.Create(dstPath)
			if err != nil {
				return err
			}
			defer dst.Close()

			if _, err := io.Copy(dst, src); err != nil {
				return err
			}

			imageName = &name
		}

		ad := domain.Ad{
			Type:         typ,
			PriceBRL:     price,
			ImagePath:    imageName,
			CEP:          formatCEP(cep),
			Street:       street,
			Number:       number,
			Complement:   complement,
			Neighborhood: neigh,
			City:         city,
			State:        state,
		}

		created, err := ads.Create(c.Context(), ad)
		if err != nil {
			return err
		}

		item := domain.AdItem{
			ID:        created.ID,
			Type:      created.Type,
			PriceBRL:  created.PriceBRL,
			CreatedAt: created.CreatedAt,
		}
		item.Address.CEP = created.CEP
		item.Address.Street = created.Street
		item.Address.Number = created.Number
		item.Address.Complement = created.Complement
		item.Address.Neighborhood = created.Neighborhood
		item.Address.City = created.City
		item.Address.State = created.State
		if created.ImagePath != nil && *created.ImagePath != "" {
			u := "/static/images/" + *created.ImagePath
			item.ImageURL = &u
		}

		return c.Status(http.StatusCreated).JSON(item)
	}
}

func ListAds(ads *service.AdsService, quotes *service.QuotesService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page := parseInt(c.Query("page"), 1)
		pageSize := parseInt(c.Query("page_size"), 10)
		if page < 1 {
			page = 1
		}
		if pageSize < 1 {
			pageSize = 10
		}
		if pageSize > 50 {
			pageSize = 50
		}

		var f repo.AdsFilter

		if v := strings.ToUpper(strings.TrimSpace(c.Query("type"))); v != "" {
			f.Type = &v
		}
		if v := strings.TrimSpace(c.Query("city")); v != "" {
			f.City = &v
		}
		if v := strings.ToUpper(strings.TrimSpace(c.Query("state"))); v != "" {
			f.State = &v
		}
		if v := strings.TrimSpace(c.Query("min_price")); v != "" {
			if x, err := strconv.ParseFloat(v, 64); err == nil {
				f.MinPrice = &x
			}
		}
		if v := strings.TrimSpace(c.Query("max_price")); v != "" {
			if x, err := strconv.ParseFloat(v, 64); err == nil {
				f.MaxPrice = &x
			}
		}

		q, err := quotes.Current(c.Context())
		if err != nil {
			return err
		}

		resp, err := ads.List(c.Context(), f, page, pageSize, q)
		if err != nil {
			return err
		}
		return c.JSON(resp)
	}
}

func parseInt(s string, def int) int {
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}

func formatCEP(cep8 string) string { 
	return cep8[:5] + "-" + cep8[5:]
}
