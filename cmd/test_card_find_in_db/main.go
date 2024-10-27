package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"main/sqlc/database"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	uid := ""
	cardUID := &uid
	flag.StringVar(cardUID, "cardUID", "", "EXAMPLE: ./<binary name> -cardUID \"<cardUID>\"")
	flag.Parse()
	uid = *cardUID
	uid = strings.ReplaceAll(uid, " ", "")
	cardUID = &uid
	if *cardUID == "" {
		flag.Usage()
		return
	}
	db, err := sql.Open(
		"mysql",
		"root:pass@tcp(localhost:3306)/kiber",
	)
	if err != nil {
		fmt.Println("Bad connection string", err)
		return
	}
	defer db.Close()
	queries := database.New(db)

	if err := db.Ping(); err != nil {
		fmt.Println("Cannot connect to database", err)
		return
	}

	// check authorization
	ctx := context.TODO()
	name, err := queries.AuthorizedCard(ctx, *cardUID)
	if err != nil {
		fmt.Println("Did not find card UID in database", err)
		return
	}
	fmt.Println("Name that was found in database based on card UID:", name)
}
