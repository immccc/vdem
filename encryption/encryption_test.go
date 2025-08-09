package encryption

import "testing"

func TestRandKeyGeneration(t *testing.T) {
	key, err := GenRandomKey()
	if err != nil {
		t.Errorf("failed to generate random key: %v", err)
	}
	if len(key) != 32 {
		t.Errorf("expected key length of 32, got %d", len(key))
	}
}

func TestEncryptionAndDecryption(t *testing.T) {
	key, err := GenRandomKey()
	if err != nil {
		t.Errorf("failed to generate random key: %v", err)
	}

	message:= "This will be encrypted for sure."
	encryptedMessage, err := Encrypt(message, key)
	if err != nil {
		t.Errorf("failed to encrypt message: %v", err)
	}

	if message == encryptedMessage {
		t.Error("expected encrypted message to differ from original message!")
	}

	decryptedMessage, err := Decrypt(encryptedMessage, key)
	if err != nil {
		t.Errorf("failed to decrypt message: %v", err)
	}

	if message != decryptedMessage {
		t.Errorf("original message: %v, decrypted msg: %v", message, decryptedMessage)
	}

}