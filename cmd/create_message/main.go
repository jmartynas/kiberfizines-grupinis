package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type request struct {
	UUID    string `json:"UUID"`
	Content string `json:"content"`
}

func main() {
	uuid := "0192c9f5-02fc-7eb1-9e72-fdf12acf481e"
	uid := ""
	message := &uid
	flag.StringVar(message, "cardUID", "", "EXAMPLE: ./<executable name> -message \"<content>\"")
	flag.Parse()
	uid = *message
	uid = strings.ReplaceAll(uid, " ", "")
	message = &uid
	if *message == "" {
		flag.Usage()
		return
	}

	publicKeyPEM, err := os.ReadFile(filepath.Join("..", "create_keys", "public.pem"))
	if err != nil {
		panic(err)
	}
	publicKeyBlock, _ := pem.Decode(publicKeyPEM)
	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		panic(err)
	}

	plaintext := []byte(*message)
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), plaintext)
	if err != nil {
		panic(err)
	}

	request := request{
		UUID:    uuid,
		Content: fmt.Sprintf("%x", ciphertext),
	}
	rq, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}

	fmt.Println("card UID input:", *message)
	fmt.Println(string(rq))
}
