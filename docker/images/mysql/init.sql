CREATE USER 'projector' IDENTIFIED BY 'projector';

CREATE DATABASE IF NOT EXISTS projection_go;
GRANT ALL ON projection_go.* TO 'projector';

USE projection_go;

CREATE TABLE subscription
(
    id VARCHAR(64) NOT NULL,
    plan_level VARCHAR(64) NOT NULL,
    plan_duration VARCHAR(64) NOT NULL,
    created DATETIME NOT NULL,
    expired DATETIME NULL,
    price NUMERIC(8,2) NOT NULL,
    status VARCHAR(64) NOT NULL,
    last_event_id VARCHAR(64) NOT NULL,
    CONSTRAINT subscription_pk PRIMARY KEY (id)
);

CREATE TABLE transaction
(
    id VARCHAR(64) NOT NULL,
    subscription_id VARCHAR(64) NOT NULL,
    created DATETIME NOT NULL,
    amount NUMERIC(8,2) NOT NULL,
    CONSTRAINT transaction_pk PRIMARY KEY (id)
);
