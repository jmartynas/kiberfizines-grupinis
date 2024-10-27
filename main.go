package main

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"main/sqlc/database"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const deviceUUID string = "0192c9f5-02fc-7eb1-9e72-fdf12acf481e"

// https://dev.to/elioenaiferrari/asymmetric-cryptography-with-golang-2ffd
func main() {
	http.HandleFunc("/", checkUser)
	if err := http.ListenAndServe(":8080", nil); err != http.ErrServerClosed {
		fmt.Println(err)
	}
}

type authorizeRequest struct {
	UUID    string `json:"UUID"`
	Content string `json:"content"`
}

func checkUser(w http.ResponseWriter, r *http.Request) {
	// read body
	body, err := io.ReadAll(r.Body)
	if errResponse(w, http.StatusBadRequest, err) {
		return
	}

	// unmarshal json request body
	// var requestContent *authorizeRequest
	requestContent := &authorizeRequest{}
	err = json.Unmarshal(body, requestContent)
	if errResponse(w, http.StatusBadRequest, err) {
		return
	}

	if requestContent.UUID != deviceUUID {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	privateKeyFile, err := os.ReadFile(filepath.Join(".", "cmd", "create_keys", "private.pem"))
	if errResponse(w, http.StatusInternalServerError, err) {
		return
	}

	privateKeyBlock, _ := pem.Decode(privateKeyFile)
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if errResponse(w, http.StatusInternalServerError, err) {
		return
	}

	encodedMessage, err := hex.DecodeString(requestContent.Content)
	if errResponse(w, http.StatusBadRequest, err) {
		return
	}

	cardUID, err := rsa.DecryptPKCS1v15(nil, privateKey, encodedMessage)
	if errResponse(w, http.StatusBadRequest, err) {
		return
	}

	cardUIDstr := strings.ReplaceAll(string(cardUID), " ", "")

	// create mysql connection
	connStr := "root:pass@tcp(localhost:3306)/kiber"
	db, err := sql.Open("mysql", connStr)
	if errResponse(w, http.StatusInternalServerError, err) {
		return
	}
	queries := database.New(db)

	if err := db.Ping(); err != nil {
		fmt.Println(err)
		return
	}

	// check authorization
	ctx := context.Background()
	name, err := queries.AuthorizedCard(ctx, cardUIDstr)
	if errResponse(w, http.StatusBadRequest, err) {
		return
	}

	if name == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fmt.Printf("%s entered\n", name)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(name))
	db.Close()
}

func errResponse(w http.ResponseWriter, status int, err error) bool {
	if err != nil && err != sql.ErrNoRows {
		fmt.Println(err)
		w.WriteHeader(status)
		_, _ = w.Write([]byte(err.Error()))
		return true
	}
	return false
}
