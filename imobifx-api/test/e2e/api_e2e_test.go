//go:build e2e

package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

var (
	baseURL string
	apiBase string
	dbDSN   string
	db      *pgxpool.Pool
	httpCli = &http.Client{Timeout: 8 * time.Second}
)

func TestMain(m *testing.M) {
	baseURL = getenv("BASE_URL", "http://localhost:8080")
	apiBase = strings.TrimRight(baseURL, "/") + "/api"
	dbDSN = getenv("TEST_DB_DSN", getenv("DB_DSN", ""))

	waitHTTP(baseURL + "/health")

	if dbDSN != "" {
		var err error
		db, err = pgxpool.New(context.Background(), dbDSN)
		if err == nil {
			_, _ = db.Exec(context.Background(), "TRUNCATE TABLE ads RESTART IDENTITY CASCADE")
			_, _ = db.Exec(context.Background(), "TRUNCATE TABLE quotes RESTART IDENTITY CASCADE")
		}
	}

	code := m.Run()

	if db != nil {
		db.Close()
	}
	os.Exit(code)
}

func TestE2E_Quotes_Create_And_Current(t *testing.T) {
	resp := doJSON(t, http.MethodPost, api("/quotes"), map[string]any{"brl_to_usd": 0.2})
	require.Equal(t, 201, resp.StatusCode)

	var created map[string]any
	readJSONInto(t, resp.Body, &created)
	require.InDelta(t, 0.2, asFloat(t, created["brl_to_usd"]), 0.000001)

	resp2 := do(t, http.MethodGet, api("/quotes/current"), nil, "")
	require.Equal(t, 200, resp2.StatusCode)

	var cur map[string]any
	readJSONInto(t, resp2.Body, &cur)
	require.InDelta(t, 0.2, asFloat(t, cur["brl_to_usd"]), 0.000001)
}

func TestE2E_Address_OK(t *testing.T) {
	resp := do(t, http.MethodGet, api("/addresses/58000000"), nil, "")
	require.Equal(t, 200, resp.StatusCode)

	var body map[string]any
	readJSONInto(t, resp.Body, &body)

	require.Equal(t, "João Pessoa", body["city"])
	require.Equal(t, "PB", body["state"])
}

func TestE2E_Address_NotFound(t *testing.T) {
	resp := do(t, http.MethodGet, api("/addresses/99999999"), nil, "")
	require.Equal(t, 404, resp.StatusCode)

	var body map[string]any
	readJSONInto(t, resp.Body, &body)
	require.Equal(t, "CEP_NOT_FOUND", body["code"])
}

func TestE2E_Ads_Create_And_List_WithUSD(t *testing.T) {
	_ = doJSON(t, http.MethodPost, api("/quotes"), map[string]any{"brl_to_usd": 0.5})

	fields := map[string]string{
		"type":         "SALE",
		"price_brl":    "100",
		"cep":          "58000-000",
		"street":       "Rua A",
		"neighborhood": "Centro",
		"city":         "João Pessoa",
		"state":        "PB",
	}

	req := newMultipartRequest(t, api("/ads"), fields, nil)
	resp := doReq(t, req)
	require.Equal(t, 201, resp.StatusCode)

	resp2 := do(t, http.MethodGet, api("/ads?page=1&page_size=10&type=SALE"), nil, "")
	require.Equal(t, 200, resp2.StatusCode)

	var body map[string]any
	readJSONInto(t, resp2.Body, &body)

	items, ok := body["items"].([]any)
	require.True(t, ok, "response must contain items[]")
	require.GreaterOrEqual(t, len(items), 1)

	first := items[0].(map[string]any)
	require.Equal(t, "SALE", first["type"])
	require.InDelta(t, 100.0, asFloat(t, first["price_brl"]), 0.0001)

	if v, ok := first["price_usd"]; ok && v != nil {
		require.InDelta(t, 50.0, asFloat(t, v), 0.0001)
	} else {
		t.Fatalf("expected price_usd in item, got keys=%v", keys(first))
	}
}

func api(path string) string {
	return apiBase + path
}

func waitHTTP(url string) {
	deadline := time.Now().Add(25 * time.Second)
	for time.Now().Before(deadline) {
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		resp, err := httpCli.Do(req)
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode >= 200 && resp.StatusCode < 500 {
				return
			}
		}
		time.Sleep(300 * time.Millisecond)
	}
}

func doJSON(t *testing.T, method, url string, payload any) *http.Response {
	t.Helper()
	b, _ := json.Marshal(payload)
	return do(t, method, url, bytes.NewReader(b), "application/json")
}

func do(t *testing.T, method, url string, body io.Reader, contentType string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(method, url, body)
	require.NoError(t, err)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	return doReq(t, req)
}

func doReq(t *testing.T, req *http.Request) *http.Response {
	t.Helper()
	resp, err := httpCli.Do(req)
	require.NoError(t, err)
	return resp
}

func readJSONInto(t *testing.T, r io.Reader, out any) {
	t.Helper()
	b, err := io.ReadAll(r)
	require.NoError(t, err)
	require.NotEmpty(t, b)
	require.NoError(t, json.Unmarshal(b, out), "body=%s", string(b))
}

func newMultipartRequest(t *testing.T, url string, fields map[string]string, file *filePart) *http.Request {
	t.Helper()

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	for k, v := range fields {
		require.NoError(t, w.WriteField(k, v))
	}

	if file != nil {
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, file.Field, file.Filename))
		h.Set("Content-Type", file.ContentType)

		part, err := w.CreatePart(h)
		require.NoError(t, err)
		_, err = part.Write(file.Bytes)
		require.NoError(t, err)
	}

	require.NoError(t, w.Close())

	req, err := http.NewRequest(http.MethodPost, url, &b)
	require.NoError(t, err)
	req.Header.Set("Content-Type", w.FormDataContentType())
	return req
}

type filePart struct {
	Field       string
	Filename    string
	ContentType string
	Bytes       []byte
}

func getenv(k, def string) string {
	if v := strings.TrimSpace(os.Getenv(k)); v != "" {
		return v
	}
	return def
}

func asFloat(t *testing.T, v any) float64 {
	t.Helper()
	f, ok := v.(float64)
	require.True(t, ok, "expected number, got %T (%v)", v, v)
	return f
}

func keys(m map[string]any) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
