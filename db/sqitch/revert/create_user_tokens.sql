-- Revert go-api-starter:create_user_tokens from pg

BEGIN;

DROP TABLE IF EXISTS user_tokens;

COMMIT;
