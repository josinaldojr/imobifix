package viacep

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/josinaldojr/imobifix-api/internal/domain"
)

var (
	ErrNotFound      = errors.New("viacep: not found")
	ErrUnavailable   = errors.New("viacep: unavailable")
	ErrInvalidAnswer = errors.New("viacep: invalid response")
)

type Client struct {
	baseURL string
	http    *http.Client
}

func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		baseURL: baseURL,
		http: &http.Client{
			Timeout: timeout,
		},
	}
}

type viaCepResp struct {
	CEP        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	UF         string `json:"uf"`
	Erro       bool   `json:"erro"`
}

func (c *Client) Lookup(ctx context.Context, cep8digits string) (domain.Address, error) {
	url := fmt.Sprintf("%s/ws/%s/json/", c.baseURL, cep8digits)

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := c.http.Do(req)
	if err != nil {
		return domain.Address{}, ErrUnavailable
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return domain.Address{}, ErrUnavailable
	}

	var v viaCepResp
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return domain.Address{}, ErrInvalidAnswer
	}
	if v.Erro {
		return domain.Address{}, ErrNotFound
	}

	return domain.Address{
		CEP:          v.CEP,
		Street:       v.Logradouro,
		Neighborhood: v.Bairro,
		City:         v.Localidade,
		State:        v.UF,
	}, nil
}
