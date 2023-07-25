-- +goose Up
-- +goose StatementBegin

CREATE SCHEMA IF NOT EXISTS slurm;

CREATE TABLE IF NOT EXISTS slurm.contact
(
    id           uuid         DEFAULT gen_random_uuid()      NOT NULL
    CONSTRAINT pk_contact
    PRIMARY KEY,
    created_at   timestamp,
    modified_at  timestamp,
    name         varchar(50)  DEFAULT '':: character varying NOT NULL,
    surname      varchar(100) DEFAULT '':: character varying NOT NULL,
    patronymic   varchar(100) DEFAULT '':: character varying NOT NULL,
    email        varchar(250),
    phone_number varchar(50),
    age          smallint
    CONSTRAINT age_check
    CHECK ((age >= 0) AND (age <= 200)),
    gender       smallint,
    is_archived  boolean      DEFAULT FALSE                  NOT NULL
    );

CREATE TABLE IF NOT EXISTS slurm."group"
(
    id            uuid      DEFAULT gen_random_uuid() NOT NULL
    CONSTRAINT pk_group
    PRIMARY KEY,
    created_at    timestamp DEFAULT CURRENT_TIMESTAMP,
    modified_at   timestamp DEFAULT CURRENT_TIMESTAMP,
    name          varchar(250)                        NOT NULL,
    description   varchar(1000),
    contact_count bigint    DEFAULT 0                 NOT NULL,
    is_archived   boolean   DEFAULT FALSE             NOT NULL
    );

CREATE TABLE IF NOT EXISTS slurm.contact_in_group
(
    id          uuid      DEFAULT gen_random_uuid() NOT NULL
    CONSTRAINT pk_contact_in_group
    PRIMARY KEY,
    created_at  timestamp DEFAULT CURRENT_TIMESTAMP,
    modified_at timestamp DEFAULT CURRENT_TIMESTAMP,
    contact_id  uuid                                not null
    constraint fk_contact_id
    references slurm.contact,
    group_id    uuid                                not null
    constraint fk_group_id
    references slurm."group"
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS slurm.contact_in_group;

DROP TABLE IF EXISTS slurm.contact;

DROP TABLE IF EXISTS slurm.group;

DROP SCHEMA IF EXISTS slurm;

-- +goose StatementEnd
