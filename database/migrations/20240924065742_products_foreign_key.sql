-- +goose Up
-- +goose StatementBegin
ALTER TABLE products
ADD CONSTRAINT FK_Admins
FOREIGN KEY (admin_id) REFERENCES admins(id)
ON DELETE CASCADE;;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE products
DROP CONSTRAINT IF EXISTS FK_Admins;
-- +goose StatementEnd
