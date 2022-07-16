
CREATE TYPE "pay" AS ENUM('cash','paid','withCard');

CREATE TABLE IF NOT EXISTS "runner"(
    "runner_id" SERIAL PRIMARY KEY NOT NULL,
    "phone_number" VARCHAR(255) NOT NULL,
    "username" VARCHAR(255) not null
);

ALTER TABLE "runner" ADD CONSTRAINT "phone_number_unique"
    UNIQUE("phone_number");
ALTER TABLE "runner" ADD CONSTRAINT "username_unique"
    UNIQUE("username");

CREATE INDEX "phone_number_idx" ON "runner"("phone_number");

CREATE TABLE IF NOT EXISTS "delivery"(
    "delivery_id" SERIAL PRIMARY KEY NOT NULL,
    "order_id" INTEGER NOT NULL,
    "runner_id" INTEGER NOT NULL,
    "reserved_at" TIMESTAMP DEFAULT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    "is_free" BOOLEAN NOT NULL DEFAULT TRUE,
    "pay" pay NOT NULL
);

ALTER TABLE "delivery" ADD CONSTRAINT "order_id_unique"
    UNIQUE("order_id");

ALTER TABLE "delivery" ADD CONSTRAINT "runner_id_fk"
    FOREIGN KEY("order_id")
    REFERENCES runner("runner_id")
    ON DELETE CASCADE;

CREATE INDEX "order_id_idx" ON "delivery" ("order_id");
CREATE INDEX "runner_id_idx" ON "delivery" ("runner_id");
