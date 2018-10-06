DO $$
BEGIN
-- Terminate all connections
IF EXISTS (SELECT 1 FROM pg_database WHERE datname = 'COST_LOGGER') THEN
    REVOKE CONNECT ON DATABASE COST_LOGGER FROM public;
    PERFORM pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = 'COST_LOGGER';
END IF;
END$$ LANGUAGE plpgsql;

-- Drop and recreate database/schema
DROP DATABASE IF EXISTS COST_LOGGER;
CREATE DATABASE COST_LOGGER WITH OWNER = 'root' ENCODING = 'UTF8';
\connect cost_logger;


-- Create tables
CREATE TABLE app_user (
    id            SERIAL       NOT NULL PRIMARY KEY,
    username      VARCHAR(100) NOT NULL,
    display_name  VARCHAR(100) NOT NULL,
    email         VARCHAR(250) NOT NULL,
    password_hash VARCHAR(100) NOT NULL
);

CREATE TABLE transaction (
    id               SERIAL       NOT NULL PRIMARY KEY,
    amount           INT          NOT NULL,
    comments         VARCHAR(500)     NULL,
    user_id          INT          NOT NULL REFERENCES app_user(id),
    last_update_date TIMESTAMP    NOT NULL DEFAULT(NOW() AT TIME ZONE 'UTC'),
    create_Date      TIMESTAMP    NOT NULL DEFAULT(NOW() AT TIME ZONE 'UTC'),
    is_active        BOOLEAN      NOT NULL DEFAULT(TRUE),
    updated_from     INT              NULL REFERENCES transaction(id)
);

CREATE TABLE transaction_user (
    transaction_id INT NOT NULL REFERENCES transaction(id),
    user_id        INT NOT NULL REFERENCES app_user(id),
    percentage     INT NOT NULL,
    PRIMARY KEY (transaction_id, user_id)
);
