package validation

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/josinaldojr/imobifix-api/internal/errors"
	"github.com/josinaldojr/imobifix-api/internal/usecase"
)

func ValidateListAdsInput(in usecase.ListAdsInput) error {
	details := fiber.Map{}

	if in.Page < 1 {
		details["page"] = "must be >= 1"
	}
	if in.PageSize < 1 || in.PageSize > 50 {
		details["page_size"] = "must be between 1 and 50"
	}
	if in.Type != nil && *in.Type != "SALE" && *in.Type != "RENT" {
		details["type"] = "must be SALE or RENT"
	}
	if in.State != nil && len(*in.State) != 2 {
		details["state"] = "must have 2 letters (UF)"
	}
	if in.MinPrice != nil && *in.MinPrice < 0 {
		details["min_price"] = "must be >= 0"
	}
	if in.MaxPrice != nil && *in.MaxPrice < 0 {
		details["max_price"] = "must be >= 0"
	}
	if in.MinPrice != nil && in.MaxPrice != nil && *in.MinPrice > *in.MaxPrice {
		details["price_range"] = "min_price must be <= max_price"
	}

	if len(details) > 0 {
		return errors.New(http.StatusBadRequest, "VALIDATION_ERROR", "Dados inválidos.", details)
	}
	
	return nil
}
