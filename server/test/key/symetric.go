package key

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"github.com/rs/zerolog/log"
)

func NewSymetric(keySize int) []byte {
	key := make([]byte, keySize)
	if _, err := rand.Read(key); err != nil {
		log.Fatal().Msgf("error while generating new symetric key: %v", err)
	}

	return key
}

func Encrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], data)

	return ciphertext, nil
}

func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	return ciphertext, nil
}

/*
func main() {
	// Key size in bytes, AES-256 requires a 32-byte key
	keySize := 32

	// Generate a random symmetric key
	key, err := generateRandomKey(keySize)
	if err != nil {
		fmt.Println("Error generating key:", err)
		return
	}

	// Convert the key to a hex string for storage or transmission
	hexKey := hex.EncodeToString(key)
	fmt.Println("Generated Symmetric Key (Hex):", hexKey)

	// Example of how to use the key for encryption and decryption
	plaintext := "Hello, World!"
	fmt.Println("Plaintext:", plaintext)

	// Encrypt
	ciphertext, err := encrypt([]byte(plaintext), key)
	if err != nil {
		fmt.Println("Error encrypting:", err)
		return
	}
	fmt.Println("Ciphertext:", ciphertext)

	// Decrypt
	decryptedPlaintext, err := decrypt(ciphertext, key)
	if err != nil {
		fmt.Println("Error decrypting:", err)
		return
	}
	fmt.Println("Decrypted Plaintext:", string(decryptedPlaintext))
}
*/
