


DROP TABLE IF EXISTS "delivery";
DROP TABLE IF EXISTS "runner";
DROP TABLE IF EXISTS "reserved";

DROP INDEX IF EXISTS "order_id_idx";
DROP INDEX IF EXISTS "runner_id_idx";
DROP INDEX IF EXISTS "phone_number_idx";

DROP TYPE IF EXISTS "pay"
