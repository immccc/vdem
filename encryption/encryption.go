package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"strings"
)
const defaultPrecision = 32
func GenRandomKey() (string, error) {

	strBuilder := strings.Builder{}
	for range defaultPrecision {
		byteOnStr, err := rand.Int(rand.Reader, big.NewInt(127 - 32)) // ASCII printable characters except DEL
		if err != nil {
			return "", fmt.Errorf("failed to generate random key: %v", err)
		}

		strBuilder.WriteByte(byte(32 + byteOnStr.Int64()))
	}

	return strBuilder.String(), nil
	

}

func Encrypt(message string, key string) (string, error) {

	aead, err := getAEAD(key)
	if err != nil {
		return "", fmt.Errorf("error on retrieving AEAD on encryption: %v", err)
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("error generating the nonce: %v", err)
	}
	msgBytes := []byte(message)
	
	encrypted := aead.Seal(nonce, nonce, msgBytes, nil)
	encryptedHex := hex.EncodeToString(encrypted)

	return encryptedHex, nil
}

func Decrypt(encryptedMessage string, key string) (string, error) {
	encrypted, err := hex.DecodeString(encryptedMessage)
	if err != nil {
		return "", fmt.Errorf("error decoding hex encrypted message: %v", err)
	}

	gcm, err := getAEAD(key)
	if err != nil {
		return "", fmt.Errorf("error on retrieving AEAD on decryption: %v", err)
	}

	decrypted, err := gcm.Open(nil, encrypted[:gcm.NonceSize()],  encrypted[gcm.NonceSize():], nil)
	if err != nil {
		return "", fmt.Errorf("error on decryption: %v", err)
	}

	decryptedStr := string(decrypted)
	return decryptedStr, nil
}

func getAEAD(key string) (cipher.AEAD, error) {
	keyBytes := []byte(key)

	if len(keyBytes) != defaultPrecision {
		return nil, fmt.Errorf("key length must be 32 bytes")
	}

	block, err := aes.NewCipher(keyBytes)

	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %v", err)
	}

	return gcm, nil
}
