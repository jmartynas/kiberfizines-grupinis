package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

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
	/*
		// create mysql connection
		ctx := context.Background()
		connStr := "root:pass@tcp(127.0.0.1:33060)/kiber"
		db, err := sql.Open("mysql", connStr)
		if errResponse(w, http.StatusBadRequest, err) {
			return
		}
		defer db.Close()
		queries := database.New(db)
	*/

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

	/*
		// get key from db
		privateKeyFile, err := queries.GetScanner(ctx, requestContent.UUID)
		if errResponse(w, http.StatusBadRequest, err) {
			return
		}
	*/
	privateKeyFile, err := os.ReadFile("./cmd/create_keys/private.pem")
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

	/*
		name, err := queries.AuthorizedCard(ctx, cardUID)
		if errResponse(w, http.StatusBadRequest, err) {
			return
		}

		log := database.InsertLogParams{
			Type: database.LogsTypeINFO,
			Message: sql.NullString{
				String: "Authorized",
				Valid:  true,
			},
			Scanner: requestContent.UUID,
			Card:    cardUID,
		}
		err = queries.InsertLog(ctx, log)
		if errResponse(w, http.StatusBadRequest, err) {
			return
		}
	*/

	// fmt.Printf("%s entered\n", name)
	// check in db if card is authorized
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(cardUID)
}

func errResponse(w http.ResponseWriter, status int, err error) bool {
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(status)
		_, _ = w.Write([]byte(err.Error()))
		return true
	}
	return false
}
