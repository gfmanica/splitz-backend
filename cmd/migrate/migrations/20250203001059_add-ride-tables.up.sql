CREATE TABLE "ride"(
    "id_ride" INTEGER NOT NULL PRIMARY KEY,
    "ds_ride" VARCHAR(255) NOT NULL,
    "vl_ride" DECIMAL(10, 2) NOT NULL,
    "qt_passengers" INTEGER NOT NULL,
    "dt_init" DATE NOT NULL,
    "dt_finish" DATE NOT NULL,
    "fg_count_weekend" BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE "person"(
    "id_person" INTEGER NOT NULL PRIMARY KEY,
    "ds_person" VARCHAR(255) NOT NULL
);

CREATE TABLE "presence"(
    "id_presence" INTEGER NOT NULL PRIMARY KEY,
    "id_person" INTEGER NOT NULL,
    "id_ride" INTEGER NOT NULL,
    "qt_presence" INTEGER NOT NULL,
    "dt_ride" DATE NOT NULL,
    CONSTRAINT "presence_id_person_foreign" FOREIGN KEY("id_person") REFERENCES "person"("id_person"),
    CONSTRAINT "presence_id_ride_foreign" FOREIGN KEY("id_ride") REFERENCES "ride"("id_ride")
);

CREATE TABLE "ride_payment"(
    "id_ride_payment" INTEGER NOT NULL PRIMARY KEY,
    "id_person" INTEGER NOT NULL,
    "id_ride" INTEGER NOT NULL,
    "vl_payment" DECIMAL(8, 2) NOT NULL,
    "dt_payment" DATE NOT NULL,
    "fg_payed" BOOLEAN NOT NULL DEFAULT FALSE,
    "ds_person" INTEGER NOT NULL,
    CONSTRAINT "ride_payment_id_person_foreign" FOREIGN KEY("id_person") REFERENCES "person"("id_person"),
    CONSTRAINT "ride_payment_id_ride_foreign" FOREIGN KEY("id_ride") REFERENCES "ride"("id_ride")
);
