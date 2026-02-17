BEGIN;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS quotes (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  brl_to_usd   NUMERIC(12,6) NOT NULL CHECK (brl_to_usd > 0),
  effective_at TIMESTAMPTZ   NOT NULL DEFAULT now(),
  created_at   TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_quotes_effective_at_desc
  ON quotes (effective_at DESC);

CREATE TABLE IF NOT EXISTS ads (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  type          TEXT NOT NULL CHECK (type IN ('SALE', 'RENT')),

  price_brl     NUMERIC(14,2) NOT NULL CHECK (price_brl >= 0),
  image_path    TEXT NULL,
  cep           TEXT NOT NULL,
  street        TEXT NOT NULL,
  number        TEXT NULL,
  complement    TEXT NULL,
  neighborhood  TEXT NOT NULL,
  city          TEXT NOT NULL,
  state         TEXT NOT NULL,

  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_ads_created_at_desc ON ads (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_ads_type ON ads (type);
CREATE INDEX IF NOT EXISTS idx_ads_city ON ads (city);
CREATE INDEX IF NOT EXISTS idx_ads_state ON ads (state);
CREATE INDEX IF NOT EXISTS idx_ads_price_brl ON ads (price_brl);

COMMIT;