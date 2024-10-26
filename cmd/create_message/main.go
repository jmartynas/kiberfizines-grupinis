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

	"github.com/google/uuid"
)

type request struct {
	UUID    string `json:"UUID"`
	Content string `json:"content"`
}

func main() {
	tmp := ""
	message := &tmp
	flag.StringVar(message, "message", "", "./<executable name> -message \"<content>\"")
	flag.Parse()

	if *message == "" {
		flag.Usage()
		return
	}

	publicKeyPEM, err := os.ReadFile("../create_keys/public.pem")
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

	uuid, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	request := request{
		UUID:    uuid.String(),
		Content: fmt.Sprintf("%x", ciphertext),
	}
	rq, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}

	fmt.Println("message input:", *message)
	fmt.Println(string(rq))
}
