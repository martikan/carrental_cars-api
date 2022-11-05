-- name: SearchCars :many
SELECT c.id
      ,c.brand_id
      ,b.name brand_name
      ,c.color
      ,c.serial
      ,c.comfort
      ,c.available
      ,c.created_at
FROM cars c
INNER JOIN brands b ON c.brand_id = b.id
WHERE (sqlc.narg(available_param)::bool IS NULL OR sqlc.narg(available_param)::bool = c.available) AND
      ((sqlc.narg(q_param)::varchar IS NULL) OR
      (LOWER(b.name) LIKE LOWER(sqlc.narg(q_param)::varchar) || '%') OR
      (LOWER(c.color) LIKE LOWER(sqlc.narg(q_param)::varchar) || '%') OR
      (LOWER(c.serial) LIKE LOWER(sqlc.narg(q_param)::varchar) || '%') OR
      (LOWER(c.comfort) LIKE LOWER(sqlc.narg(q_param)::varchar) || '%'))
LIMIT sqlc.arg(limit_param)::int
OFFSET sqlc.arg(offset_param)::int;

-- name: GetAllCars :many
SELECT c.id
      ,c.brand_id
      ,b.name brand_name
      ,c.color
      ,c.serial
      ,c.comfort
      ,c.available
      ,c.created_at
FROM cars c
INNER JOIN brands b ON c.brand_id = b.id
LIMIT $1
OFFSET $2;

-- name: GetCarById :one
SELECT c.id
      ,c.brand_id
      ,b.name brand_name
      ,c.color
      ,c.serial
      ,c.comfort
      ,c.available
      ,c.created_at
FROM cars c
INNER JOIN brands b ON c.brand_id = b.id
WHERE c.id = $1
LIMIT 1;

-- name: CreateCar :one
INSERT INTO cars (brand_id, color, "serial", comfort)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateCar :one
UPDATE cars
SET brand_id = $2,
color = $3,
"serial" = $4,
comfort = $5,
available = $6
WHERE id = $1
RETURNING *;

-- name: DeleteCar :exec
DELETE FROM cars WHERE id = $1;
