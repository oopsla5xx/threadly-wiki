-- name: GetIAMUserByEmail :one
SELECT *
FROM iam_users
WHERE LOWER(email) = LOWER($1)
    AND deleted_at IS NULL;

-- name: UpdateIAMUserStatus :one
UPDATE iam_users
SET status     = $2,
    updated_at = NOW()
WHERE id = $1
  AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteIAMUser :execrows
UPDATE iam_users
SET deleted_at = NOW(),
    updated_at = NOW()
WHERE id = $1
  AND deleted_at IS NULL;
