-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS variants (
    id SERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuid_generate_v4(),
    variant_name VARCHAR(30) NOT NULL,
    quantity INT,
    product_id INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION setTimeStamp_Variants()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_timestamp_variants
BEFORE UPDATE ON variants
FOR EACH ROW
EXECUTE FUNCTION setTimeStamp_Variants();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS set_timestamp_variants ON variants;
DROP FUNCTION IF EXISTS setTimeStamp_Variants();
DROP TABLE IF EXISTS variants;
-- +goose StatementEnd