package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"main/sqlc/database"
	"net/http"
	"strings"
	"time"

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
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", checkUser)
	http.HandleFunc("/logs", logs)

	fmt.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != http.ErrServerClosed {
		fmt.Println("Server error:", err)
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

	cardUID = Unpad(cardUID)
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
		errResponse(w, http.StatusServiceUnavailable, fmt.Errorf("database unavailable"))
		return
	}

	// check authorization
	ctx := context.Background()
	name, err := queries.AuthorizedCard(ctx, cardUIDstr)
	if err == sql.ErrNoRows {
		fmt.Printf("Access denied for card: %s\n", cardUIDstr)
		log := database.InsertLogParams{
			Uid:       cardUIDstr,
			Permitted: false,
			Time:      time.Now().Truncate(time.Second),
		}
		_ = queries.InsertLog(ctx, log)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if errResponse(w, http.StatusBadRequest, err) {
		return
	}

	if name == "" {
		fmt.Printf("Access denied: no name found for card %s\n", cardUIDstr)
		log := database.InsertLogParams{
			Uid:       cardUIDstr,
			Permitted: false,
			Time:      time.Now().Truncate(time.Second),
		}
		_ = queries.InsertLog(ctx, log)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	log := database.InsertLogParams{
		Uid:       cardUIDstr,
		Permitted: true,
		Time:      time.Now().Truncate(time.Second),
	}
	if err := queries.InsertLog(ctx, log); err != nil {
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

func Unpad(data []byte) []byte {
	n := int(data[len(data)-1])
	return data[:len(data)-n]
}

func logs(w http.ResponseWriter, r *http.Request) {
	connStr := "root:pass@tcp(127.0.0.1:3306)/kiber?parseTime=true"
	db, err := sql.Open("mysql", connStr)
	if errResponse(w, http.StatusInternalServerError, err) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("failed to create database connection: %v\n", err)
		return
	}
	queries := database.New(db)

	if err := db.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("failed to database: %v\n", err)
		return
	}

	// check authorization
	ctx := context.Background()
	logsRows, err := queries.SelectLogs(ctx)
	if errResponse(w, http.StatusBadRequest, err) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("failed to get logs from database: %v\n", err)
		return
	}

	type tmp struct {
		ID        int32
		Uid       string
		Permitted bool
		Time      time.Time
		UserName  string
	}

	var logs []tmp
	for _, v := range logsRows {
		log := tmp{
			ID:        v.ID,
			Uid:       v.Uid,
			Permitted: v.Permitted,
			Time:      v.Time,
			UserName:  v.UserName.String,
		}
		logs = append(logs, log)
	}

	tmpl := template.Must(template.ParseFiles("./templates/index.html"))
	if err := tmpl.Execute(w, logs); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("failed to create and send html: %v\n", err)
	}
}
