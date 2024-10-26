-- name: AuthorizedCard :one
SELECT user_name FROM card
WHERE uid = ? LIMIT 1;
