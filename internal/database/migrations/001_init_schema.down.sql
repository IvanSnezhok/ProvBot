-- Drop indexes
DROP INDEX IF EXISTS idx_admin_users_telegram_id;
DROP INDEX IF EXISTS idx_outages_location;
DROP INDEX IF EXISTS idx_outages_status;
DROP INDEX IF EXISTS idx_bot_logs_created_at;
DROP INDEX IF EXISTS idx_bot_logs_level;
DROP INDEX IF EXISTS idx_message_logs_created_at;
DROP INDEX IF EXISTS idx_message_logs_telegram_id;
DROP INDEX IF EXISTS idx_message_logs_user_id;
DROP INDEX IF EXISTS idx_users_telegram_id;

-- Drop tables
DROP TABLE IF EXISTS admin_users;
DROP TABLE IF EXISTS outages;
DROP TABLE IF EXISTS bot_logs;
DROP TABLE IF EXISTS message_logs;
DROP TABLE IF EXISTS users;

