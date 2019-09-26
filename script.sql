create table if not exists engines
(
    engine_id     serial  not null
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

create table if not exists transmissions
(
    transmission_model        text    not null,
    transmission_type         text    not null,
    transmission_gears_number integer not null,
    transmission_id           serial  not null
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
    brand_id              serial not null
        constraint brands_pk
            primary key,
    brand_name            text   not null,
    brand_creator_country text   not null
);

alter table brands
    owner to kolya59;

create table if not exists cars
(
    model           text    not null,
    brand_id        integer not null
        constraint cars___brands
            references brands
            on update cascade on delete cascade,
    engine_id       integer not null
        constraint cars___engines
            references engines
            on update cascade on delete cascade,
    transmission_id integer not null
        constraint cars___transmissions
            references transmissions
            on update cascade on delete cascade,
    price           integer not null,
    constraint cars_pk
        primary key (model, brand_id, engine_id, transmission_id)
);

alter table cars
    owner to kolya59;

create unique index if not exists brands_brand_id_uindex
    on brands (brand_id);


