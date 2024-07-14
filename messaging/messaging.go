package messaging

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/decred/dcrd/dcrec/secp256k1/v4/schnorr"
)

type MessageType string

const (
	TagTypeUser = "p"
)

type Message struct {
	Id        string   `json:"id"`
	PubKey    string   `json:"pubKey"`
	CreatedAt uint64   `json:"createdAt"`
	Kind      uint16   `json:"kind"`
	Tags      [][]string `json:"tags"`
	Content   string   `json:"content"`
	Sig       string   `json:"sig"`
}

func (msg *Message) Build() string {
	js, _ := json.MarshalIndent(msg, " ", " ")
	return string(js)
}

func (msg *Message) Sign(privateKey string) {
	hashId := sha256.New()

	
	id := fmt.Sprintf("[0,%v,%v,%v,%v]", msg.CreatedAt, msg.Kind, msg.Tags, msg.Content)
	hashIdAsBytes := []byte(id)
	hashId.Write(hashIdAsBytes)

	msg.Id = hex.EncodeToString(hashId.Sum(nil))

	privateKeyAsHex := hex.EncodeToString([]byte(privateKey))
	privkeyBytes, _ := hex.DecodeString(privateKeyAsHex)
	privKey := secp256k1.PrivKeyFromBytes(privkeyBytes)

	storedHashIdBytes, _ := hex.DecodeString(msg.Id)
	signature, _ := schnorr.Sign(privKey, storedHashIdBytes)
	msg.Sig = hex.EncodeToString(signature.Serialize())

	pubKeyBytes := privKey.PubKey().SerializeCompressed()

	msg.PubKey = hex.EncodeToString(pubKeyBytes)

}

func IsSigned(msg *Message) bool {
	return msg.Id != "" && msg.Sig != ""
}

func (msg *Message) Verify() bool {
	pubKeyBytes, _ := hex.DecodeString(msg.PubKey)
	pubKey, _ := secp256k1.ParsePubKey(pubKeyBytes)

	sigBytes, _ := hex.DecodeString(msg.Sig)
	signature, _ := schnorr.ParseSignature(sigBytes)

	idBytes, _ := hex.DecodeString(msg.Id)
	return signature.Verify(idBytes[:], pubKey)
}
