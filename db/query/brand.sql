-- name: GetBrandById :one
SELECT id
      ,"name"
FROM brands
WHERE id = $1
LIMIT 1;

-- name: CountBrandByName :one
SELECT COUNT(id)
FROM brands
WHERE LOWER("name") = LOWER($1)
LIMIT 1;

-- name: CreateBrand :one
INSERT INTO brands ("name")
VALUES ($1)
RETURNING *;

-- name: DeleteBrand :exec
DELETE FROM brands WHERE id = $1;
