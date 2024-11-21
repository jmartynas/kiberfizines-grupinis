-- name: AuthorizedCard :one
SELECT user_name FROM card
WHERE uid = ? LIMIT 1;

-- name: InsertLog :exec
INSERT INTO log (
uid,
permitted,
time
) values (
?, ?, ?
);

-- name: SelectLogs :many
Select * from log
LEFT JOIN card on card.uid = log.uid;
