-- Grant permissions to database user
-- IMPORTANT: Replace 'provbot' with your actual database user name from .env (POSTGRES_USER)
-- Usage: psql -d provbot_db -U postgres -f 002_grant_permissions.sql
-- Or connect as superuser and run: psql -d provbot_db -c "GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO provbot;"

-- Grant privileges on existing tables
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO provbot;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO provbot;

-- Grant default privileges for future tables/sequences
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO provbot;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO provbot;

