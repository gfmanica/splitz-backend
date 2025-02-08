CREATE TABLE "ride"(
    "id_ride" SERIAL PRIMARY KEY,
    "ds_ride" VARCHAR(255) NOT NULL,
    "vl_ride" DECIMAL(10, 2) NOT NULL,
    "dt_init" DATE NOT NULL,
    "dt_finish" DATE NOT NULL,
    "fg_count_weekend" BOOLEAN NOT NULL DEFAULT FALSE
);


CREATE TABLE "ride_payment"(
    "id_ride_payment" SERIAL PRIMARY KEY,
    "id_ride" INTEGER NOT NULL,
    "vl_payment" DECIMAL(8, 2) NOT NULL,
    "fg_payed" BOOLEAN NOT NULL DEFAULT FALSE,
    "ds_person" VARCHAR(255) NOT NULL,
    CONSTRAINT "ride_payment_id_ride_foreign" FOREIGN KEY("id_ride") REFERENCES "ride"("id_ride")
);

CREATE TABLE "presence"(
    "id_presence" SERIAL PRIMARY KEY,
    "id_ride_payment" INTEGER NOT NULL,
    "qt_presence" INTEGER NOT NULL,
    "dt_ride" DATE NOT NULL,
    CONSTRAINT "presence_id_ride_payment_foreign" FOREIGN KEY("id_ride_payment") REFERENCES "ride_payment"("id_ride_payment")
);
