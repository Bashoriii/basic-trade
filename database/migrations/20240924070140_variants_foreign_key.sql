-- +goose Up
-- +goose StatementBegin
ALTER TABLE variants
ADD CONSTRAINT FK_Products
FOREIGN KEY (product_id) REFERENCES products(id)
ON DELETE CASCADE;;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE variants
DROP CONSTRAINT IF EXISTS FK_Products;
-- +goose StatementEnd
