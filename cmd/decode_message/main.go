package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"strings"
)

type test struct {
	UUID    string `json:"UUID"`
	Content string `json:"content"`
}

func main() {
	key := []byte{
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

	iv := []byte{
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

	body := flag.String("body", "", "http request body")
	flag.Parse()
	*body = strings.ReplaceAll(*body, " ", "")
	if *body == "" {
		flag.Usage()
		return
	}

	test := &test{}
	if err := json.Unmarshal([]byte(*body), test); err != nil {
		fmt.Println(err)
		return
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
		return
	}

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

func Unpad(data []byte, blockSize int) []byte {
	n := int(data[len(data)-1])
	return data[:len(data)-n]
}
