package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"main/sqlc/database"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const deviceUUID string = "0192c9f5-02fc-7eb1-9e72-fdf12acf481e"

var key = []byte{
	0x2b,
	0x7e,
	0x15,
	0x16,
	0x28,
	0xae,
	0xd2,
	0xa6,
	0xab,
	0xf7,
	0x97,
	0x99,
	0x89,
	0xcf,
	0xab,
	0x12,
}

// https://dev.to/elioenaiferrari/asymmetric-cryptography-with-golang-2ffd
func main() {
	fmt.Println("Server started on :8080")
	http.HandleFunc("/", checkUser)
	if err := http.ListenAndServe(":8080", nil); err != http.ErrServerClosed {
		fmt.Println("Server error:", err) // Modified error message
	}
}

type authorizeRequest struct {
	UUID    string `json:"UUID"`
	IV      string `json:"iv"`
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
		fmt.Println("Unauthorized device attempt")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	block, err := aes.NewCipher(key)
	if errResponse(w, http.StatusInternalServerError, err) {
		return
	}

	iv, err := hex.DecodeString(requestContent.IV)
	if errResponse(w, http.StatusBadRequest, err) {
		return
	}

	encodedMessage, err := hex.DecodeString(requestContent.Content)
	if errResponse(w, http.StatusBadRequest, err) {
		return
	}

	cardUID := make([]byte, len(encodedMessage))
	dec := cipher.NewCBCDecrypter(block, iv)
	dec.CryptBlocks(cardUID, encodedMessage)

	fmt.Printf("Received UID bytes: %v\n", cardUID)

	cardUID = Unpad(cardUID, aes.BlockSize)
	cardUIDstr := strings.ToUpper(hex.EncodeToString(cardUID))

	// create mysql connection
	connStr := "root:pass@tcp(127.0.0.1:3306)/kiber"
	db, err := sql.Open("mysql", connStr)
	if errResponse(w, http.StatusInternalServerError, err) {
		return
	}

	defer db.Close()

	queries := database.New(db)

	if err := db.Ping(); err != nil {
		fmt.Println("Database connection failed:", err)
		return
	}

	// check authorization
	ctx := context.Background()
	name, err := queries.AuthorizedCard(ctx, cardUIDstr)
	if err == sql.ErrNoRows {
		fmt.Printf("Access denied for card: %s\n", cardUIDstr)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if errResponse(w, http.StatusBadRequest, err) {
		return
	}

	if name == "" {
		fmt.Printf("Access denied: no name found for card %s\n", cardUIDstr)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fmt.Printf("Access granted: %s entered (%s)\n", name, cardUIDstr)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(name))
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

func Unpad(data []byte, blockSize int) []byte {
	n := int(data[len(data)-1])
	return data[:len(data)-n]
}

/*

	ciphertext, err := hex.DecodeString((*test).Content)
	if err != nil {
		fmt.Println(err)
		return
	}

	plaintext := make([]byte, len(ciphertext))
	dec := cipher.NewCBCDecrypter(block, iv)
	dec.CryptBlocks(plaintext, ciphertext)

	plaintext = Unpad(plaintext, aes.BlockSize)

	fmt.Println("Decrypted UID:", string(plaintext))
}

*/
