package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

func TrimIndex(text []byte) int {
	for i:=0; i<len(text); i++ {
		if text[i] == 0 {
			return i
		}
	}

	return len(text)
}

func Encrypt(payload, secret string) (cipherString string, err error) {
	var key [32]byte;
	copy(key[:], []byte(secret))

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
	return string(text[:TrimIndex(text)]), nil
}

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
	
	return string(plaintext[:TrimIndex(plaintext)]), nil
}
