-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "Messages" (
      "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
      "text" TEXT NOT NULL,
      "is_processed" BOOLEAN NOT NULL DEFAULT false,
      "processed_at" TIMESTAMP NULL,
      "created_at" TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "Messages";
-- +goose StatementEnd
