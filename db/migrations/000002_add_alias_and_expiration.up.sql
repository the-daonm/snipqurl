-- Alter short_code length and add expires_at column
ALTER TABLE urls ALTER COLUMN short_code TYPE varchar(30);
ALTER TABLE urls ADD COLUMN IF NOT EXISTS expires_at timestamp;
