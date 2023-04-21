-- +migrate Up
ALTER TABLE "users" ADD "is_deleted_account" bool NOT NULL DEFAULT false;

-- +migrate Down
ALTER TABLE "users" DROP COLUMN "is_deleted_account";