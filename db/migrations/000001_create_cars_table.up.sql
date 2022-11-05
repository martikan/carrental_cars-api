CREATE TABLE IF NOT EXISTS "brands" (
    "id" bigserial PRIMARY KEY,
    "name" varchar(255) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS "cars" (
    "id" bigserial PRIMARY KEY,
    "brand_id" bigint NOT NULL,
    "color" varchar(100) NOT NULL,
    "serial" varchar(20) NOT NULL,
    "comfort" varchar(1) NOT NULL CHECK (comfort in ('S', 'A', 'B', 'C', 'D')),
    "available" boolean NOT NULL DEFAULT TRUE,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    foreign key (brand_id) references "brands" (id)
);

CREATE INDEX "fk_cars_brand_idx" ON "cars" (brand_id);
CREATE INDEX "cars_search_terms" ON "cars" (color, comfort, available, "serial");