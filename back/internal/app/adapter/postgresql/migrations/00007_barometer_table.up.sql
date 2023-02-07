CREATE TABLE "barometers" (
    id          bigserial
        primary key,
    created_at  timestamp with time zone,
    updated_at  timestamp with time zone,
    deleted_at  timestamp with time zone,
    "date"      timestamp without time zone not null,
    "user_id"   bigint not null,
    "score"     int not null,
    "type"      text not null,
    CONSTRAINT "fk_barometers_user" FOREIGN KEY ("user_id") REFERENCES "users"("id")
);

CREATE UNIQUE INDEX "uniq_barometer_per_day_type_user" ON "barometers" ("user_id", "type", "date");
