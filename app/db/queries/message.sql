-- name: CreateMessage :one
INSERT INTO
  messages (
    message_type,
    room_id,
    client_id,
    created_at,
    data
  )
VALUES
  ($1, $2, $3, $4, $5)
RETURNING
  id,
  message_type,
  room_id,
  client_id,
  created_at,
  data;

-- name: GetMessagesByRoomID :many
SELECT
  id,
  message_type,
  room_id,
  client_id,
  created_at,
  data
FROM
  messages
WHERE
  room_id = $1
ORDER BY
  created_at;

-- name: GetPreviousMessages :many
SELECT
  id,
  message_type,
  room_id,
  client_id,
  created_at,
  data
FROM
  messages
WHERE
  room_id = $1
  AND created_at < $2
ORDER BY
  created_at
LIMIT
  $3;