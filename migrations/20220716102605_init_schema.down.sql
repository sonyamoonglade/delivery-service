


DROP TABLE IF EXISTS "reserved";
DROP TABLE IF EXISTS "telegram_runner";
DROP TABLE IF EXISTS "runner";
DROP TABLE IF EXISTS "delivery";

DROP INDEX IF EXISTS "order_id_idx";
DROP INDEX IF EXISTS "runner_id_idx";
DROP INDEX IF EXISTS "phone_number_idx";

DROP TYPE IF EXISTS "pay"
