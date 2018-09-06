DO $$
BEGIN
-- Terminate all connections
IF EXISTS (SELECT 1 FROM pg_database WHERE datname = 'cost_logger') THEN
    REVOKE CONNECT ON DATABASE cost_logger FROM public;
    PERFORM pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = 'cost_logger';
END IF;
END$$ LANGUAGE plpgsql;

-- Drop and recreate database/schema
DROP DATABASE IF EXISTS cost_logger;
CREATE DATABASE cost_logger WITH OWNER = 'root' ENCODING = 'UTF8';
\connect cost_logger;
CREATE SCHEMA cost_logger;


-- Create tables
CREATE TABLE cost_logger.USER (
    Id       SERIAL,
    Name     VARCHAR(100) NOT NULL,
    Password VARCHAR(100) NOT NULL
);

CREATE TABLE cost_logger.TRANSACTION (
    Id             SERIAL,
    Amount         DECIMAL   NOT NULL,
    Date           TIMESTAMP NOT NULL,
    UserId         INT       NOT NULL REFERENCES USER(Id),
    LastUpdateTime TIMESTAMP NOT NULL
);

CREATE TABLE cost_logger.TRANSACTION_USER (
    TransactionId INT NOT NULL REFERENCES TRANSACTION(Id),
    UserId        INT NOT NULL REFERENCES USER(Id),
    Percentage    DECIMAL NOT NULL
);
