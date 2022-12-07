package main

import (
	"fmt"

	"github.com/Murchik/kuznyechik/crypt"
)

func main() {
	key1 := []byte("someFancyKey1")
	key2 := []byte("someFancyKey2")

	crypt.ExpandKey(key1, key2)
	fmt.Printf("key 1: '%s'\nkey 2: '%s'\n\n", key1, key2)

	message := []byte("Hello world! ")
	fmt.Printf("Before encryption:\nlen = %v\n'%s'\n\n", len(message), message)

	encryptMsg := crypt.Encrypt(message)
	fmt.Printf("After encryption:\nlen = %v\ndata = %v\n\n", len(encryptMsg), encryptMsg)

	decryptMsg := crypt.Decrypt(encryptMsg)
	fmt.Printf("After decryption:\nlen = %v\n'%s'\n", len(decryptMsg), decryptMsg)
}
