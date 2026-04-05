DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'manifest_verifications') THEN
    RAISE EXCEPTION 'missing table: manifest_verifications';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'release_verifications') THEN
    RAISE EXCEPTION 'missing table: release_verifications';
  END IF;
  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_name = 'manifest_verifications' AND column_name = 'details'
  ) THEN
    RAISE EXCEPTION 'missing column: manifest_verifications.details';
  END IF;
END $$;
