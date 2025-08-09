package vote

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"immccc/vdem/messaging"
	"immccc/vdem/peer"
	"math/big"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	uuid "github.com/google/uuid"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
)

type PollStatus int

const (
	Open PollStatus = 0
	Closed
)

type Poll struct {
	Id           uuid.UUID
	Description  string
	Options      []string
	Status       PollStatus
	AllowedPeers []string
}

type Vote struct {
	Height         			uint64
	Comment                 string
	unencodedSelection      uint
	EncodedSelection		string
	UpdatedOptions          []string
	signature	            schnorr.Signature
	swapsDone			 	uint8
}

func CreateVote(optionIdx uint) (*Vote, error) {
	/// Creates a vote for the given option index. The vote is created with a randomly generated key

	vote := &Vote{
		unencodedSelection: optionIdx,
		Comment:   "TODO Comments not implemented yet",
	}

	voteKey, err := genVoteKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate vote key: %v", err)
	}

	err = signVote(vote, voteKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign vote: %v", err)
	}

	return vote, nil
}


// Shuffle the vote with a random peer of the network. This step is necessary to ensure that the vote is not directly linked to the voter.
func RequestShuffleVoteWithRandomPeer(vote *Vote, peers[] *peer.Peer) error {

	idxPeer, err := rand.Int(rand.Reader, big.NewInt(int64(len(peers))))
    if err != nil {
		return fmt.Errorf("failed to get arandom peer: %v", err)
    }

	peer := peers[idxPeer.Int64()]

	event := messaging.SwapVotesEvent(peer.PubKey, vote.EncodedSelection, vote.EncodedSelection)
	peer.SendMessage(messaging.BuildEventMessage(&event))

	return nil
	
}

func genVoteKey() (string, error) {
	privKey, err := btcec.NewPrivateKey()
	if err != nil {
		return "", err
	}

	return privKey.Key.String(), nil
}

// Sign the vote with the private key of the voter. The signature is used to verify the authenticity of the vote.
// The signature is generated using the private key of the voter and the option that was voted for.
// The signature is then stored in the vote object.
func signVote(vote *Vote, privateKey string) error {

	// TODO Probably a common method, to avoid repetition to what's present in event.go?
	privateKeyAsHex := hex.EncodeToString([]byte(privateKey))
	privkeyBytes, _ := hex.DecodeString(privateKeyAsHex)
	privKey := secp256k1.PrivKeyFromBytes(privkeyBytes)

	optionIdxBytes, _ := hex.DecodeString(fmt.Sprintf("%v", vote.unencodedSelection))
	signature, error := schnorr.Sign(privKey, optionIdxBytes)
	if error != nil {
		return fmt.Errorf("failed to sign vote: %v", error)
	}

	vote.signature = *signature
	
	sigBytes := signature.Serialize()
	vote.EncodedSelection = hex.EncodeToString(sigBytes)

	return nil
}