DO $$
BEGIN
IF NOT EXISTS (SELECT column_name FROM information_schema.columns WHERE table_name='transaction' and column_name='updated_from') THEN
    ALTER TABLE transaction ADD COLUMN updated_from INT NULL REFERENCES transaction(id);
END IF;
END
$$;
