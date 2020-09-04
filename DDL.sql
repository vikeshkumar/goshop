-- CREATE ROLE goshop WITH
--     LOGIN
--     NOSUPERUSER
--     NOCREATEDB
--     NOCREATEROLE
--     INHERIT
--     NOREPLICATION
--     CONNECTION LIMIT -1
--     PASSWORD 'xxxxxx';
-- COMMENT ON ROLE goshop IS 'goshop';
--
-- CREATE DATABASE goshop
--     WITH
--     OWNER = goshop
--     ENCODING = 'UTF8'
--     CONNECTION LIMIT = -1;
CREATE TABLE IF NOT EXISTS SSO_PROVIDER
(
    name          varchar(100) PRIMARY KEY,
    client_id     varchar(255) NOT NULL,
    creation_time timestamp DEFAULT CURRENT_TIMESTAMP,
    modified_time timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ROLE
(
    code          varchar(20) PRIMARY KEY,
    description   varchar(100) NOT NULL,
    expired       bool      DEFAULT FALSE,
    locked        bool      DEFAULT FALSE,
    creation_time timestamp DEFAULT CURRENT_TIMESTAMP,
    modified_time timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS APP_USER
(
    uid           varchar(20) PRIMARY KEY,
    email         varchar(50) NOT NULL,
    locked        bool      default false,
    hash          varchar(20),
    password      varchar(20),
    local_user    bool      DEFAULT TRUE,
    sso_user      bool      DEFAULT FALSE,
    sso_provider  SSO_PROVIDER,
    creation_time timestamp DEFAULT CURRENT_TIMESTAMP,
    modified_time timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS USER_ROLES
(
    role          varchar(20),
    uid           varchar(20),
    foreign key (role) references ROLE (code),
    foreign key (uid) references APP_USER (uid),
    creation_time timestamp DEFAULT CURRENT_TIMESTAMP,
    modified_time timestamp DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS WEBSITE
(
    uid           varchar(200) PRIMARY KEY,
    title         varchar(200) NOT NULL,
    pattern       varchar(200) ARRAY,
    creation_time timestamp DEFAULT CURRENT_TIMESTAMP,
    modified_time timestamp DEFAULT CURRENT_TIMESTAMP
);
