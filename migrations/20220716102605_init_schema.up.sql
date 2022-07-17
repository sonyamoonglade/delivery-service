

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
    "created_at" TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
    "is_free" BOOLEAN NOT NULL DEFAULT TRUE,
    "pay" pay NOT NULL
);

ALTER TABLE "delivery" ADD CONSTRAINT "order_id_unique"
    UNIQUE("order_id");

CREATE INDEX "order_id_idx" ON "delivery" ("order_id");


CREATE TABLE IF NOT EXISTS "reserved"(
    "delivery_id" SERIAL PRIMARY KEY NOT NULL,
    "runner_id" INTEGER NOT NULL,
    "reserved_at" TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc')
);

ALTER TABLE "reserved" ADD CONSTRAINT "runner_id_fk"
    FOREIGN KEY("runner_id")
    REFERENCES runner("runner_id")
    ON DELETE CASCADE;

ALTER TABLE "reserved" ADD CONSTRAINT "delivery_id_fk"
    FOREIGN KEY("delivery_id")
    REFERENCES delivery("delivery_id")
    ON DELETE CASCADE;

CREATE INDEX "runner_id_idx" ON "reserved" ("runner_id");

CREATE TABLE IF NOT EXISTS "telegram_runner"(
    "runner_id" INTEGER PRIMARY KEY NOT NULL,
    "telegram_id" INTEGER NOT NULL
);

ALTER TABLE "telegram_runner" ADD CONSTRAINT "runner_id_fk_tg"
    FOREIGN KEY("runner_id")
    REFERENCES runner("runner_id")
    ON DELETE CASCADE;

ALTER TABLE "telegram_runner" ADD CONSTRAINT "runner_tg_id_unique"
    UNIQUE (runner_id, telegram_id);

