-- Revert short_code length and remove expires_at column
-- Note: This might fail if any short_code is longer than 10 characters
ALTER TABLE urls ALTER COLUMN short_code TYPE varchar(10);
ALTER TABLE urls DROP COLUMN expires_at;
