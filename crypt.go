package main

/*
	Written by Resyfer
	https://github.com/resyfer/pswdmngr/blob/main/crypt.go
*/

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

/**

	3 Hours of Pain and Suffering later, finally I completed coding these
	functions to meet my needs

**/


func trimIndex(text []byte) int {
	for i:=0; i<len(text); i++ {
		if text[i] == 0 {
			return i
		}
	}

	return len(text)
}

func checkNewLine(text []byte) bool {

	for i:=0; i<len(text); i++ {
		if text[i] == byte(10) || text[i] == byte(32) {
			return true
		}
	}
	return false
}

// Encrypt(String, Key) gives the gibberish string and an error.
// Make sure Key is smaller than 32 bytes (len(Key) <= 32)
func Encrypt(payload, secret string) (cipherString string, err error) {
	var key [32]byte;
	copy(key[:], []byte(secret))

	encryptedByteString := []byte{10}
	for checkNewLine(encryptedByteString) {

		cipherText, err := aes.NewCipher(key[:])
		if err != nil {
			return "", err
		}
		
		gcm, err := cipher.NewGCM(cipherText)
		if err != nil {
			return "", err
		}
	
		nonce := make([]byte, gcm.NonceSize())
		if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
			return "", err
		}
	
		text := gcm.Seal(nonce, nonce, []byte(payload), nil)
		
		encryptedByteString = text[:trimIndex(text)]
	}

	return string(encryptedByteString), nil
}

// Decrypt(encryptedString, key) gives back original string and error.
// Make sure Key is smaller than 32 bytes (len(Key) <= 32)
func Decrypt(cipherString, secret string) (payload string , err error) {
	var key [32]byte;
	copy(key[:], []byte(secret))
	
	ciphertext := []byte(cipherString)

	c, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("too small")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("bad")
	}
	
	return string(plaintext[:trimIndex(plaintext)]), nil
}
