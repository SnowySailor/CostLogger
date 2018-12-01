DO $$
BEGIN
IF NOT EXISTS (SELECT column_name FROM information_schema.columns WHERE table_name='transaction' and column_name='updated_from') THEN
    ALTER TABLE transaction ADD COLUMN updated_from INT NULL REFERENCES transaction(id);
END IF;
END
$$;

DO $$
BEGIN
IF NOT EXISTS (SELECT column_name FROM information_schema.columns WHERE table_name='transaction_user' AND column_name='is_paid') THEN
    ALTER TABLE transaction_user ADD COLUMN is_paid BOOLEAN NULL;
END IF;
UPDATE transaction_user SET is_paid = false WHERE is_paid is null;
ALTER TABLE transaction_user ALTER COLUMN is_paid SET NOT NULL;
END
$$;