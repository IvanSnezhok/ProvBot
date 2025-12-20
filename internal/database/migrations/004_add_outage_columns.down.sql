-- Remove group_id and street columns from outages table
DROP INDEX IF EXISTS idx_outages_group_id;
ALTER TABLE outages DROP COLUMN IF EXISTS street;
ALTER TABLE outages DROP COLUMN IF EXISTS group_id;
