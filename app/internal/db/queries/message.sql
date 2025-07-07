-- name: CreateMessage :one
INSERT INTO
  messages (id, room_id, client_id, created_at, read, content)
VALUES
  ($1, $2, $3, $4, $5, $6)
RETURNING
  *;

-- name: BulkInsertMessages :copyfrom
INSERT INTO
  messages (id, room_id, client_id, created_at, read, content)
VALUES
  ($1, $2, $3, $4, $5, $6);

-- name: GetAllMessagesByRoomID :many
SELECT
  *
FROM
  messages
WHERE
  room_id = $1
ORDER BY
  created_at;

-- name: GetAllMessagesBeforeCreatedAt :many
SELECT
  *
FROM
  messages
WHERE
  room_id = $1
  AND created_at < $2
ORDER BY
  created_at
LIMIT
  $3;