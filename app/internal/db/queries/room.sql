-- name: CreateRoom :one
INSERT INTO
  rooms (item_id, client1, client2)
VALUES
  ($1, $2, $3)
RETURNING
  *;

-- name: GetRoomById :one
SELECT
  *
FROM
  rooms
WHERE
  id = $1;

-- name: GetAllRoomsByClient :many
SELECT
  *
FROM
  rooms
WHERE
  $1::uuid IN (client1, client2);

-- name: GetRoomByItemAndClients :one
SELECT
  *
FROM
  rooms
WHERE
  item_id = $1
  AND client1 = $2
  AND client2 = $3;

-- name: GetRoomByIdAndClient :one
SELECT
  *
FROM
  rooms
WHERE
  id = $1
  AND @client_id IN (client1, client2);

-- name: DeleteRoomById :exec
DELETE FROM rooms
WHERE
  id = $1;