# ImobiFX

Sistema com API em Go e frontend Flutter desktop para cadastro e listagem de anuncios imobiliarios, com cotacao BRL -> USD e consulta de endereco por CEP.

## Estrutura

- `imobifx-api`: backend REST em Go
- `imobifx-frontend`: frontend Flutter desktop

## Funcionalidades

- Cadastro de anuncios (`SALE` e `RENT`)
- Upload opcional de imagem do imovel
- Consulta de CEP via backend (integracao ViaCEP)
- Fallback para preenchimento manual do endereco quando CEP falha
- Cadastro de cotacoes BRL -> USD
- Listagem paginada de anuncios com filtros
- Exibicao de preco em BRL e USD
- Internacionalizacao no frontend (PT e EN via parametro)
- Documentacao Swagger/OpenAPI da API

## Stack Tecnica

- Backend: Go 1.24, Fiber v2, PostgreSQL 16, Docker/Compose, ViaCEP, OpenAPI/Swagger
- Frontend: Flutter desktop, Dart, `http`, `file_picker`, `intl`

## Decisoes Tecnicas Adotadas

- Separacao clara de responsabilidades no backend (handlers, service, repo, validation, domain), para manter manutenibilidade e testes mais previsiveis.
- Persistencia em PostgreSQL por consistencia transacional e bom suporte a filtros/paginacao para listagem de anuncios.
- API REST com validacao no backend como fonte da verdade, mesmo com validacao no frontend, para garantir seguranca e consistencia.
- Consulta de CEP encapsulada no backend para isolar a dependencia externa (ViaCEP) e padronizar tratamento de falhas para o frontend.
- Upload de imagem no endpoint de anuncio via `multipart/form-data`, mantendo os metadados e arquivo no mesmo fluxo de criacao.
- Conversao de valores BRL -> USD usando cotacao vigente persistida, evitando regra de negocio no frontend.
- Frontend com i18n por `--dart-define=APP_LOCALE`, facilitando execucao em idioma diferente sem alterar codigo.
- Placeholder local para imagem ausente/quebrada, evitando UI degradada.
- Swagger/OpenAPI exposto pela API para facilitar inspecao, teste manual e integracao cliente-servidor.

## Requisitos

- Docker + Docker Compose
- Flutter SDK 3.22+ (com suporte desktop habilitado, ex: Windows)

## Como Rodar

1. Subir backend + banco (na raiz do projeto):

```bash
docker compose up --build
```

2. API disponivel em:

- `http://localhost:8080`
- Swagger UI: `http://localhost:8080/swagger`
- OpenAPI YAML: `http://localhost:8080/swagger/openapi.yaml`

3. Rodar frontend (em outro terminal):

```bash
cd imobifx-frontend
flutter pub get
flutter run -d windows --dart-define=API_BASE_URL=http://localhost:8080 --dart-define=APP_LOCALE=pt
```

4. Rodar frontend em ingles:

```bash
flutter run -d windows --dart-define=API_BASE_URL=http://localhost:8080 --dart-define=APP_LOCALE=en
```

## Comandos Uteis

- Subir tudo: `docker compose up --build ou make up`
- Derrubar ambiente: `docker compose down -v ou make down`
- Logs da API: `cd imobifx-api && make logs`
- Testes backend: `cd imobifx-api && make test`
- Testes e2e backend: `cd imobifx-api && make e2e`
