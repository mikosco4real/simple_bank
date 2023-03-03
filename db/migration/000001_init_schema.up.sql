-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS users_id_seq;

-- Table Definition
CREATE TABLE "public"."users" (
    "id" int8 NOT NULL DEFAULT nextval('users_id_seq'::regclass),
    "uuid" varchar NOT NULL,
    "firstname" varchar NOT NULL,
    "lastname" varchar NOT NULL,
    "email" varchar NOT NULL,
    "password" varchar NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY ("id")
);

-- Indexes
CREATE INDEX "users_firstname" ON "public"."users" USING BTREE ("firstname");
CREATE INDEX "users_lastname" ON "public"."users" USING BTREE ("lastname");
CREATE UNIQUE INDEX "users_email" ON "public"."users" USING BTREE ("email");
CREATE INDEX "users_password" ON "public"."users" USING BTREE ("password");
CREATE INDEX "users_name" ON "public"."users" USING BTREE ("firstname","lastname");
CREATE INDEX "users_auth" ON "public"."users" USING BTREE ("email","password");

-- Column Comment
COMMENT ON COLUMN "public"."users"."uuid" IS 'uuid is generated automatically in the code';