-- +migrate Up
ALTER TABLE "equipment"
	ALTER COLUMN "inventory_number" TYPE varchar(255) USING "inventory_number"::varchar;

-- +migrate Down
ALTER TABLE "equipment"
	ALTER COLUMN "inventory_number" TYPE bigint USING
		NULLIF(regexp_replace("inventory_number", '[^0-9]', '', 'g'), '')::bigint;
