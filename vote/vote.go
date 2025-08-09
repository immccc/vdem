package vote

import (
	// "encoding/hex"
	"fmt"
	"immccc/vdem/encryption"
	"immccc/vdem/peer"
	"maps"

	// "github.com/decred/dcrd/dcrec/secp256k1/v4"
	uuid "github.com/google/uuid"

	"github.com/btcsuite/btcd/btcec/v2"
)

type PollStatus int

const (
	Open PollStatus = 0
	Closed
)

type Poll struct {
	Id                         uuid.UUID
	Description                string
	Options                    []string
	HoldKey                    string
	Status                     PollStatus
	AllowedPeers               map[string]int
	peersWithUnshuffledVote    map[string]int
	localVote                  *Vote
	ongoingShufflingWithPeerId *string
}

type Vote struct {
	Comment            string
	unencodedSelection uint
	EncodedSelection   string
	UpdatedOptions     []string
}

func (poll *Poll) SetPeers(peers ...*peer.Peer) {
	if poll.AllowedPeers == nil {
		poll.AllowedPeers = make(map[string]int)
	}

	for _, peer := range peers {
		poll.AllowedPeers[peer.PubKey] = 1
	}

	poll.peersWithUnshuffledVote = maps.Clone(poll.AllowedPeers)

}

func (poll *Poll) CreateVote(optionIdx uint) error {
	/// Creates a vote for the given option index. The vote is created with a randomly generated key
	if poll.localVote != nil {
		return fmt.Errorf("vote already created for this poll")
	}

	vote := &Vote{
		unencodedSelection: optionIdx,
		Comment:            "TODO Comments not implemented yet",
	}

	voteKey, err := genVoteKey()
	if err != nil {
		return fmt.Errorf("failed to generate vote key: %v", err)
	}

	vote.EncodedSelection, err = encryption.Encrypt(fmt.Sprintf("%d", vote.unencodedSelection), voteKey)
	if err != nil {
		return fmt.Errorf("failed to sign vote: %v", err)
	}

	poll.HoldKey = voteKey
	poll.localVote = vote

	return nil
}

func (poll *Poll) GetLocalVote() *Vote {
	return poll.localVote
}

func (poll *Poll) GetOnGoingShufflingPeer() *string {
	return poll.ongoingShufflingWithPeerId
}

func (poll *Poll) ResetOngoingshufflingPeer() {
	poll.ongoingShufflingWithPeerId = nil
}

func (poll *Poll) SetOngoingShufflingPeer(fromPeerKey string) {
	poll.ongoingShufflingWithPeerId = &fromPeerKey
}

func (poll *Poll) GetPeersWithUnshuffledVote() map[string]int {
	if poll.peersWithUnshuffledVote == nil {
		poll.peersWithUnshuffledVote = make(map[string]int)
	}

	return poll.peersWithUnshuffledVote
}

func genVoteKey() (string, error) {
	privKey, err := btcec.NewPrivateKey()
	if err != nil {
		return "", err
	}

	return privKey.Key.String(), nil
}
