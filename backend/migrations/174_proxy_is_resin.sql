-- proxies: Resin mode expands auth username to Platform.{account_id} at request time
ALTER TABLE proxies ADD COLUMN IF NOT EXISTS is_resin BOOLEAN NOT NULL DEFAULT FALSE;
