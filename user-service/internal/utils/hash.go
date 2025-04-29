package utils

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func Hash512(orderId, statusCode, grossAmount string) (string, error) {

	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	if serverKey == "" {
		return "", errors.New("env MIDTRANS_SERVER_KEY is not set")
	}

	input := orderId + statusCode + grossAmount + serverKey
	inputBytes := []byte(input)
	sha512Hasher := sha512.New()
	sha512Hasher.Write(inputBytes)
	hashedInputBytes := sha512Hasher.Sum(nil)
	hashedInputString := hex.EncodeToString(hashedInputBytes)

	return hashedInputString, nil
}

func HashPassword(password string) (string, error) {

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func ComparePassword(hashedPassword, password string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
