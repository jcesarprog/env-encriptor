package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
)

func main() {
	key := "mysecretkey12345" // Replace with a strong, secret key
	inputFile := ".env"
	outputFile := "env.enc"

	//* Encrypting

	err := encryptFile(inputFile, outputFile, key)
	if err != nil {
		panic(err)
	}

	// Clean up the original file
	err = os.Remove(inputFile)
	if err != nil {
		fmt.Println("There was an error deleting the file, please remove ir manually.")
		return
	}

	// * Decripting

	// Restore the original .env
	err = decryptFile(outputFile, inputFile, key)
	if err != nil {
		panic(err)
	}

	fmt.Println("File decrypted successfully.")

}

// encryptFile encrypts the data from the inputFile and writes it to the outputFile.
func encryptFile(inputFile, outputFile, key string) error {
	plaintext, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	err = os.WriteFile(outputFile, ciphertext, 0644)
	return err
}

func decryptFile(inputFile, outputFile, key string) error {
	ciphertext, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return err
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCTR(block, iv)
	plaintext := make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)

	err = os.WriteFile(outputFile, plaintext, 0644)
	return err
}
