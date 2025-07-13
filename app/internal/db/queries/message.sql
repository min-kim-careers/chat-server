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

-- name: GetMessagesBeforeCreatedAt :many
SELECT
  *
FROM
  (
    SELECT
      *
    FROM
      messages
    WHERE
      room_id = $1
      AND created_at < $2
    ORDER BY
      created_at DESC
    LIMIT
      $3
  ) AS latest
ORDER BY
  created_at ASC;
