CREATE ROLE goshop WITH
    LOGIN
    NOSUPERUSER
    NOCREATEDB
    NOCREATEROLE
    INHERIT
    NOREPLICATION
    CONNECTION LIMIT -1
    PASSWORD 'xxxxxx';
COMMENT ON ROLE goshop IS 'goshop';

CREATE DATABASE goshop
    WITH
    OWNER = goshop
    ENCODING = 'UTF8'
    CONNECTION LIMIT = -1;