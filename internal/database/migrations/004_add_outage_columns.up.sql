-- Add group_id and street columns to outages table
ALTER TABLE outages ADD COLUMN IF NOT EXISTS group_id VARCHAR(255);
ALTER TABLE outages ADD COLUMN IF NOT EXISTS street VARCHAR(255);

-- Create index for faster outage lookups by group
CREATE INDEX IF NOT EXISTS idx_outages_group_id ON outages(group_id) WHERE status = 'active';
