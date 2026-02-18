package swaggerdocs

import _ "embed"

var (
	//go:embed openapi.yaml
	openAPI []byte
)

func OpenAPIYAML() []byte {
	return openAPI
}
