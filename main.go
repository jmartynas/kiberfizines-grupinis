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

	_ "github.com/go-sql-driver/mysql" // Import MySQL driver
)

const deviceUUID string = "0192c9f5-02fc-7eb1-9e72-fdf12acf481e" // UUID for the device to authorize

// AES key for encryption and decryption
var key = []byte{
	0x2b, 0x7e, 0x15, 0x16, 0x28, 0xae, 0xd2, 0xa6,
	0xab, 0xf7, 0x97, 0x99, 0x89, 0xcf, 0xab, 0x12,
}

func main() {
	// Serve static files (e.g., images, stylesheets)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Define routes
	http.HandleFunc("/", checkUser) // Route for checking user authorization
	http.HandleFunc("/logs", logs)  // Route to view logs

	fmt.Println("Server started on :8080")

	// Start HTTP server on port 8080
	if err := http.ListenAndServe(":8080", nil); err != http.ErrServerClosed {
		fmt.Println("Server error:", err)
	}
}

// Struct to parse the incoming authorization request
type authorizeRequest struct {
	UUID    string `json:"UUID"`
	IV      string `json:"iv"`
	Content string `json:"content"`
}

// Handler to check the user (authorization process)
func checkUser(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if errResponse(w, http.StatusBadRequest, err) {
		return
	}

	// Unmarshal the JSON request body into the authorizeRequest struct
	requestContent := &authorizeRequest{}
	err = json.Unmarshal(body, requestContent)
	if errResponse(w, http.StatusBadRequest, err) {
		return
	}

	// Check if the UUID matches the authorized device UUID
	if requestContent.UUID != deviceUUID {
		fmt.Println("Unauthorized device attempt")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Initialize AES cipher block using the pre-defined key
	block, err := aes.NewCipher(key)
	if errResponse(w, http.StatusInternalServerError, err) {
		return
	}

	// Decode the IV (initialization vector) and the encrypted content from the request
	iv, err := hex.DecodeString(requestContent.IV)
	if errResponse(w, http.StatusBadRequest, err) {
		return
	}

	encodedMessage, err := hex.DecodeString(requestContent.Content)
	if errResponse(w, http.StatusBadRequest, err) {
		return
	}

	// Decrypt the content using CBC mode and AES cipher
	cardUID := make([]byte, len(encodedMessage))
	dec := cipher.NewCBCDecrypter(block, iv)
	dec.CryptBlocks(cardUID, encodedMessage)

	// Log the received UID bytes
	fmt.Printf("Received UID bytes: %v\n", cardUID)

	// Unpad the decrypted data
	cardUID = Unpad(cardUID)

	// Convert the UID to a hex string and make it uppercase
	cardUIDstr := strings.ToUpper(hex.EncodeToString(cardUID))

	// Create MySQL connection string
	connStr := "root:pass@tcp(127.0.0.1:3306)/kiber"
	db, err := sql.Open("mysql", connStr)
	if errResponse(w, http.StatusInternalServerError, err) {
		return
	}
	defer db.Close()

	// Create queries instance from database package
	queries := database.New(db)

	// Check if the database is available
	if err := db.Ping(); err != nil {
		errResponse(w, http.StatusServiceUnavailable, fmt.Errorf("database unavailable"))
		return
	}

	// Check if the card UID is authorized
	ctx := context.Background()
	name, err := queries.AuthorizedCard(ctx, cardUIDstr)
	if err == sql.ErrNoRows {
		// If card UID is not authorized, log the attempt and return unauthorized
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

	// If the card has no associated name, deny access
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

	// Log the successful authorization and grant access
	log := database.InsertLogParams{
		Uid:       cardUIDstr,
		Permitted: true,
		Time:      time.Now().Truncate(time.Second),
	}
	if err := queries.InsertLog(ctx, log); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Print a success message and respond with the name associated with the card UID
	fmt.Printf("Access granted: %s entered (%s)\n", name, cardUIDstr)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(name))
}

// Helper function to handle error responses
func errResponse(w http.ResponseWriter, status int, err error) bool {
	if err != nil && err != sql.ErrNoRows {
		fmt.Println(err)
		w.WriteHeader(status)
		_, _ = w.Write([]byte(err.Error()))
		return true
	}
	return false
}

// Unpad the decrypted data (for AES CBC mode padding)
func Unpad(data []byte) []byte {
	n := int(data[len(data)-1]) // Last byte contains the padding length
	return data[:len(data)-n]   // Remove padding
}

// Handler to view logs
func logs(w http.ResponseWriter, r *http.Request) {
	// Create MySQL connection string for logs
	connStr := "root:pass@tcp(127.0.0.1:3306)/kiber?parseTime=true"
	db, err := sql.Open("mysql", connStr)
	if errResponse(w, http.StatusInternalServerError, err) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("failed to create database connection: %v\n", err)
		return
	}
	queries := database.New(db)

	// Check if the database is available
	if err := db.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("failed to database: %v\n", err)
		return
	}

	// Retrieve logs from the database
	ctx := context.Background()
	logsRows, err := queries.SelectLogs(ctx)
	if errResponse(w, http.StatusBadRequest, err) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("failed to get logs from database: %v\n", err)
		return
	}

	// Create a temporary struct to hold log data for display
	type tmp struct {
		ID        int32
		Uid       string
		Permitted bool
		Time      time.Time
		UserName  string
	}

	// Populate the logs slice with the retrieved rows
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

	// Parse and execute the HTML template to render the logs
	tmpl := template.Must(template.ParseFiles("./templates/index.html"))
	if err := tmpl.Execute(w, logs); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("failed to create and send html: %v\n", err)
	}
}
