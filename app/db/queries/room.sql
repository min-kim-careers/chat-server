-- name: CreateRoom :one
INSERT INTO
  rooms (slug, item_id, buyer_id, seller_id)
VALUES
  ($1, $2, $3, $4)
RETURNING
  id,
  slug,
  item_id,
  buyer_id,
  seller_id,
  created_at,
  updated_at;

-- name: GetRoomBySlug :one
SELECT
  id,
  slug,
  item_id,
  buyer_id,
  seller_id,
  created_at,
  updated_at
FROM
  rooms
WHERE
  slug = $1;
