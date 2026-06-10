-- name: CreateUser :one
-- CreateUser inserts a new user record into the database.
-- It returns the full row of the newly created user, including automatic timestamps and UUIDs.
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetUser :one
-- GetUser retrieves a single user's details by their unique name.
-- This is primarily used during login to verify an account exists.
SELECT * FROM users
where name = $1;

-- name: DeleteUsers :exec
-- DeleteUsers purges all users from the database.
-- This is a destructive operation used primarily by the developer 'reset' command.
DELETE FROM users;

-- name: GetUsers :many
-- GetUsers retrieves a list of all registered usernames.
-- This is used to output a list of users currently registered with the application.
SELECT name from users;