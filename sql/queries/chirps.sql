-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, user_id, message)
VALUES (
  gen_random_uuid(), now(), now(), $1, $2
)
RETURNING *;

-- name: ListChirps :many
SELECT * FROM chirps
ORDER BY created_at ASC;

