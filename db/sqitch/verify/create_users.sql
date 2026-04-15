-- Verify go-fiber-template:create_users on pg

SELECT id, email, name, hashed_password, created_at, updated_at
  FROM users
 WHERE FALSE;
