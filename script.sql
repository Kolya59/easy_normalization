create table if not exists engines
(
    engine_id     uuid  not null
        constraint engines_pk
            primary key,
    engine_model  text    not null,
    engine_power  integer not null,
    engine_volume integer,
    engine_type   text    not null
);

alter table engines
    owner to kolya59;

create unique index if not exists engines_engine_id_uindex
    on engines (engine_id);

create unique index if not exists engines_engine_model_uindex
    on engines (engine_model);

create table if not exists transmissions
(
    transmission_model        text    not null,
    transmission_type         text    not null,
    transmission_gears_number integer not null,
    transmission_id           uuid  not null
        constraint transmissions_pk
            primary key
);

alter table transmissions
    owner to kolya59;

create unique index if not exists transmissions_transmission_model_uindex
    on transmissions (transmission_model);

create unique index if not exists transmissions_transmission_id_uindex
    on transmissions (transmission_id);

create table if not exists brands
(
    brand_id              uuid not null
        constraint brands_pk
            primary key,
    brand_name            text   not null,
    brand_creator_country text   not null
);

alter table brands
    owner to kolya59;

create unique index if not exists brands_brand_id_uindex
    on brands (brand_id);

create unique index if not exists brands_brand_name_uindex
    on brands (brand_name);

create table if not exists wheels
(
    wheel_id     uuid not null
        constraint wheels_pk
            primary key,
    wheel_radius integer not null,
    wheel_color  text    not null,
    wheel_model  text    not null
);

alter table wheels
    owner to kolya59;

create unique index if not exists wheels_wheel_id_uindex
    on wheels (wheel_id);

create unique index if not exists wheels_wheel_model_uindex
    on wheels (wheel_model);

create table if not exists cars
(
    id uuid not null primary key,
    model           text    not null,
    brand_id        uuid not null
        constraint cars___brands
            references brands,
    engine_id       uuid not null
        constraint cars___engines
            references engines,
    transmission_id uuid not null
        constraint cars___transmissions
            references transmissions,
    price           integer not null,
    wheel_id        uuid not null
        constraint cars___wheels
            references wheels
);

alter table cars
    owner to kolya59;

