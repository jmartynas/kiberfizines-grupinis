package main

import (
	"bytes"
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
	IV      string `json:"iv"`
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

	text := flag.String("cardUID", "", "card UID")
	flag.Parse()

	*text = strings.ReplaceAll(*text, " ", "")

	if *text == "" {
		flag.Usage()
		return
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
		return
	}

	plain := Pad([]byte(*text), aes.BlockSize)
	ciphertext := make([]byte, len(plain))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plain)

	test := test{
		UUID:    "0192c9f5-02fc-7eb1-9e72-fdf12acf481e",
		IV:      hex.EncodeToString(iv),
		Content: fmt.Sprintf("%x", ciphertext),
	}
	body, err := json.Marshal(test)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
}

func Pad(data []byte, blockSize int) []byte {
	n := blockSize - len(data)%blockSize
	padding := bytes.Repeat([]byte{byte(n)}, n)
	return append(data, padding...)
}
