package event

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/decred/dcrd/dcrec/secp256k1/v4/schnorr"
)

type EventType string

const (
	TagTypeUser = "p"
)

type Event struct {
	Id        string     `json:"id"`
	PubKey    string     `json:"pubKey"`
	CreatedAt uint64     `json:"createdAt"`
	Kind      uint16     `json:"kind"`
	Tags      [][]string `json:"tags"`
	Content   string     `json:"content"`
	Sig       string     `json:"sig"`
}

func (event *Event) Build() string {
	js, _ := json.MarshalIndent(event, " ", " ")
	return string(js)
}

func (event *Event) Sign(privateKey string) {
	hashId := sha256.New()

	id := fmt.Sprintf("[0,%v,%v,%v,%v]", event.CreatedAt, event.Kind, event.Tags, event.Content)
	hashIdAsBytes := []byte(id)
	hashId.Write(hashIdAsBytes)

	event.Id = hex.EncodeToString(hashId.Sum(nil))

	privateKeyAsHex := hex.EncodeToString([]byte(privateKey))
	privkeyBytes, _ := hex.DecodeString(privateKeyAsHex)
	privKey := secp256k1.PrivKeyFromBytes(privkeyBytes)

	storedHashIdBytes, _ := hex.DecodeString(event.Id)
	signature, _ := schnorr.Sign(privKey, storedHashIdBytes)
	event.Sig = hex.EncodeToString(signature.Serialize())

	pubKeyBytes := privKey.PubKey().SerializeCompressed()

	event.PubKey = hex.EncodeToString(pubKeyBytes)

}

func IsSigned(event *Event) bool {
	return event.Id != "" && event.Sig != ""
}

func (event *Event) Verify() bool {
	pubKeyBytes, _ := hex.DecodeString(event.PubKey)
	pubKey, _ := secp256k1.ParsePubKey(pubKeyBytes)

	sigBytes, _ := hex.DecodeString(event.Sig)
	signature, _ := schnorr.ParseSignature(sigBytes)

	idBytes, _ := hex.DecodeString(event.Id)
	return signature.Verify(idBytes[:], pubKey)
}
