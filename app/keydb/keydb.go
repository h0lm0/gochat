package keydb

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	dbmodels "gochat/database/models"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

var Client *redis.Client

func InitClient() {
	keydbPassword := os.Getenv("KEYDB_PASSWORD")
	keydbHost := "keydb"

	Client = redis.NewClient(&redis.Options{
		Addr:     keydbHost + ":6379",
		Password: keydbPassword,
		DB:       0,
	})

	_, err := Client.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Error connecting to keydb: %v", err))
	}

	log.Printf("Successfully connected to keydb")
}

func SetKey(key string, value interface{}) error {
	err := Client.Set(ctx, key, value, 0).Err()
	if err != nil {
		log.Printf("Error setting key to keydb: %v", err)
		return err
	}
	return nil
}

func GetKey(key string) (string, error) {
	value, err := Client.Get(ctx, key).Result()
	return value, err
}

func AddEncryptionKey(user dbmodels.User) (bool, error) {
	key := user.Username + "-rsa"
	check, _ := GetKey(key + "-pub")

	if check != "" {
		log.Println("Encryption key for user " + user.Username + " already exists")
		return true, nil
	}

	log.Println("Creating encryption key for user " + user.Username)
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Println("Error generating keypair for user " + user.Username + ": " + err.Error())
		return false, err
	}

	pub := privateKey.Public()

	keyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)

	pubPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(pub.(*rsa.PublicKey)),
		},
	)

	privErr := SetKey(key+"-priv", keyPEM)
	if privErr != nil {
		log.Println(privErr.Error())
	}

	pubErr := SetKey(key+"-pub", pubPEM)
	if pubErr != nil {
		log.Println(pubErr.Error())
	}

	return true, nil

}

func GetPublicKey(user dbmodels.User) (*rsa.PublicKey, error) {
	key := user.Username + "-rsa"

	pubKey, _ := GetKey(key + "-pub")

	publicKeyBlock, _ := pem.Decode([]byte(pubKey))
	publicKey, err := x509.ParsePKCS1PublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return publicKey, nil

}

func GetPrivateKey(user dbmodels.User) (*rsa.PrivateKey, error) {
	key := user.Username + "-rsa"

	privKey, _ := GetKey(key + "-priv")

	privateKeyBlock, _ := pem.Decode([]byte(privKey))
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		panic(err)
	}

	return privateKey, nil

}
