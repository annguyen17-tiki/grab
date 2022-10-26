CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create table if not exists "accounts" (
    "id" uuid not null default uuid_generate_v4(),
    "username" varchar(256) not null,
    "password" varchar(256) not null,
    "firstname" varchar(256) not null,
    "lastname" varchar(256) not null,
    "phone" varchar(16) not null,
    "avatar" text,
    "role" varchar(16) not null,
    "created_at" timestamp not null,
    "updated_at" timestamp not null,
    "deleted_at" timestamp,

    primary key ("id"),
    constraint "unq_account_username" unique ("username"),
    constraint "unq_account_phone" unique ("phone")
);

create table if not exists "locations" (
    "account_id" uuid not null,
    "longitude" double precision not null,
    "latitude" double precision not null,
    "geo_hash" varchar(16) not null,
    "created_at" timestamp not null,
    "updated_at" timestamp not null,
    "deleted_at" timestamp,

    primary key ("account_id"),
    constraint "fk_locations_account_id" foreign key ("account_id") references "accounts"("id")
);

create table if not exists "bookings" (
    "id" uuid not null default uuid_generate_v4(),
    "user_id" uuid not null,
    "driver_id" uuid,
    "vehicle" varchar(16),
    "from_longitude" double precision not null,
    "from_latitude" double precision not null,
    "to_longitude" double precision not null,
    "to_latitude" double precision not null,
    "status" varchar(16) not null,
    "created_at" timestamp not null,
    "updated_at" timestamp not null,
    "deleted_at" timestamp,

    primary key ("id"),
    constraint "fk_bookings_user_id" foreign key ("user_id") references "accounts"("id"),
    constraint "fk_bookings_driver_id" foreign key ("user_id") references "accounts"("id")
);

create table if not exists "drivers" (
    "account_id" uuid not null,
    "vehicle" varchar(16) not null,
    "status" varchar(16) not null,
    "created_at" timestamp not null,
    "updated_at" timestamp not null,
    "deleted_at" timestamp,

    primary key ("account_id"),
    constraint "fk_drivers_account_id" foreign key ("account_id") references "accounts"("id")
);

create table if not exists "offers" (
    "booking_id" uuid not null,
    "driver_id" uuid not null,
    "status" varchar(16) not null,
    "created_at" timestamp not null,
    "updated_at" timestamp not null,
    "deleted_at" timestamp,

    primary key ("booking_id", "driver_id"),
    constraint "fk_offers_booking_id" foreign key ("booking_id") references "bookings"("id"),
    constraint "fk_offers_driver_id" foreign key ("driver_id") references "drivers"("account_id")
);

create table if not exists "notifications" (
    "id" uuid not null default uuid_generate_v4(),
    "account_id" uuid not null,
    "status" varchar(16) not null,
    "content" json not null,

    "created_at" timestamp not null,
    "updated_at" timestamp not null,
    "deleted_at" timestamp,

    primary key ("id"),
    constraint "fk_notifications_account_id" foreign key ("account_id") references "accounts"("id")
);