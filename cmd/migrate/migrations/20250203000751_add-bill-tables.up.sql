CREATE TABLE "bill"(
    "id_bill" SERIAL PRIMARY KEY,
    "ds_bill" VARCHAR(255) NOT NULL,
    "vl_bill" DECIMAL(10, 2) NOT NULL,
    "qt_person" INTEGER NOT NULL
);

CREATE TABLE "bill_payment"(
    "id_payment" SERIAL PRIMARY KEY,
    "vl_payment" DECIMAL(10, 2) NOT NULL,
    "dt_payment" DATE,
    "ds_person" VARCHAR(255) NOT NULL,
    "fg_payed" BOOLEAN NOT NULL DEFAULT FALSE,
    "fg_custom_payment" BOOLEAN NOT NULL DEFAULT FALSE,
    "id_bill" INTEGER NOT NULL,
    CONSTRAINT "bill_payment_id_bill_foreign" FOREIGN KEY("id_bill") REFERENCES "bill"("id_bill")
);
