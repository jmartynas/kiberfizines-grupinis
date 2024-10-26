-- name: GetScanner :one
SELECT private_key FROM scanner
WHERE uuid = ? LIMIT 1;

-- name: AuthorizedCard :one
SELECT user_name FROM card
WHERE uid = ? LIMIT 1;

-- name: InsertLog :exec
INSERT INTO logs 
(type, message, scanner, card)
VALUES
(?, ?, ?, ?);
