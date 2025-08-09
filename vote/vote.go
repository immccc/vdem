package vote

import (
	"encoding/hex"
	"fmt"
	"immccc/vdem/peer"

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
	Id 	  	     	uuid.UUID
	Description		string
	Options      	[]string
	Status       	PollStatus
	AllowedPeers 	map[string]int
	PeersWithVote 	map[string]int
	localVote    	*Vote
}

type Vote struct {
	Comment                 string
	unencodedSelection      uint
	EncodedSelection		string
	UpdatedOptions          []string
	signature	            schnorr.Signature
	swapsDone			 	uint8
}

func (p* Poll) SetPeers(peers ...*peer.Peer)  {
	if p.AllowedPeers != nil {
		p.AllowedPeers = make(map[string]int)
	}

	for _, peer := range peers {
		p.AllowedPeers[peer.PubKey] = 1
	}

}

func (poll* Poll) CreateVote(optionIdx uint) (error) {
	/// Creates a vote for the given option index. The vote is created with a randomly generated key
	if poll.localVote != nil {
		return fmt.Errorf("vote already created for this poll")
	}


	vote := &Vote{
		unencodedSelection: optionIdx,
		Comment:   "TODO Comments not implemented yet",
	}

	voteKey, err := genVoteKey()
	if err != nil {
		return fmt.Errorf("failed to generate vote key: %v", err)
	}

	err = signVote(vote, voteKey)
	if err != nil {
		return fmt.Errorf("failed to sign vote: %v", err)
	}

	poll.localVote = vote

	return nil
}

func (poll *Poll) GetLocalVote() (*Vote) {
	return poll.localVote
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