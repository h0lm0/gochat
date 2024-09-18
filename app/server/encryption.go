package server

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	dbmodels "gochat/database/models"
	"gochat/keydb"
)

func EncryptData(user dbmodels.User, data string) (string, error) {
	if len(data) > 250 {
		return "", errors.New("data too long for encryption")
	}
	publicKey, pubKeyErr := keydb.GetPublicKey(user)
	if pubKeyErr != nil {
		panic(pubKeyErr)
	}
	plaintext := []byte(data)
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plaintext)
	if err != nil {
		return "", err
	}
	return string(ciphertext), nil
}

func DecryptData(user dbmodels.User, data []byte) (string, error) {
	privateKey, err := keydb.GetPrivateKey(user)
	if err != nil {
		panic(err)
	}

	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, data)
	if err != nil {
		panic(err)
	}

	return string(plaintext), nil
}
