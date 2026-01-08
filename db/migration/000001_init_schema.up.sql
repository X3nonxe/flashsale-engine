CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_hash" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "products" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "description" text,
  "image_url" varchar,
  "price" decimal(10, 2) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "product_stocks" (
  "id" bigserial PRIMARY KEY,
  "product_id" bigint NOT NULL,
  "quantity" int NOT NULL DEFAULT 0,
  "version" int NOT NULL DEFAULT 1,
  "last_updated" timestamptz NOT NULL DEFAULT (now()),
  CONSTRAINT "fk_product" FOREIGN KEY ("product_id") REFERENCES "products" ("id"),
  CONSTRAINT "stock_not_negative" CHECK (quantity >= 0)
);

CREATE TABLE "orders" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "user_id" bigint NOT NULL,
  "product_id" bigint NOT NULL,
  "quantity" int NOT NULL,
  "status" varchar NOT NULL DEFAULT 'PENDING',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  CONSTRAINT "fk_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id"),
  CONSTRAINT "fk_order_product" FOREIGN KEY ("product_id") REFERENCES "products" ("id")
);

CREATE INDEX "idx_product_stocks_product_id" ON "product_stocks" ("product_id");
CREATE INDEX "idx_orders_user_id" ON "orders" ("user_id");
CREATE UNIQUE INDEX "idx_unique_order" ON "orders" ("user_id", "product_id");
