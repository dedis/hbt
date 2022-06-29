// package state implements the state handling of our transaction family. Note
// that the serialization / deserialization must be deterministic, according to
// Sawtooth documentation:
//
//   https://sawtooth.hyperledger.org/docs/1.2/architecture/global_state.html
//
// > It is critical to select a serialization scheme which is deterministic
// > across executions of the transaction, across platforms, and across
// > versions of the serialization framework.

package state

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"strings"

	"ledger/financial/serde"

	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/suites"
)

// KeyAddress defines the key address of the financial smart contract state
var KeyAddress = hexdigest("hbt_financial")[:70]

var suite = suites.MustFind("ed25519")

// WalletID is a hex-encoded public key - kyber point
type WalletID string

// ToPoint returns the point representation of the walletID
func (w WalletID) ToPoint() (kyber.Point, error) {
	buff, err := hex.DecodeString(string(w))
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex: %v", err)
	}

	p := suite.Point()

	err = p.UnmarshalBinary(buff)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal point: %v", err)
	}

	return p, nil
}

// newState returns a new initialized state
func newState() state {
	return state{
		Wallets: Wallets{},
	}
}

// state defines the data stored by the smart contract
type state struct {
	Wallets Wallets
}

// NewHandler returns a new initialized Handler
func NewHandler(ctx *processor.Context, serde serde.Serde) (Handler, error) {
	state, err := loadState(ctx, serde)
	if err != nil {
		return Handler{}, fmt.Errorf("failed to load state: %v", err)
	}

	return Handler{
		context: ctx,
		serde:   serde,
		state:   state,
	}, nil
}

// Handler handles addressing and (de)serialization of the state
type Handler struct {
	context *processor.Context
	serde   serde.Serde
	state   state
}

// CreateWallet creates a new wallet with walletID. Returns an error if the
// wallet already exists.
func (h Handler) CreateWallet(wid WalletID) error {
	_, found := h.state.Wallets.Get(wid)
	if found {
		return fmt.Errorf("wallet with id %q already exists", wid)
	}

	h.state.Wallets.Set(Wallet{
		Balance: 0,
		State:   WalletOK,
	})

	err := h.StoreState()
	if err != nil {
		return fmt.Errorf("failed to store state: %v", err)
	}

	return nil
}

// GetWallet returns the wallet corresponding to the walletID, or an error if
// not found.
func (h Handler) GetWallet(wid WalletID) (Wallet, error) {
	wallet, found := h.state.Wallets.Get(wid)
	if !found {
		return Wallet{}, fmt.Errorf("wallet with id %q already exists", wid)
	}

	return wallet, nil
}

// StoreState saves the current state
func (h Handler) StoreState() error {
	buff, err := h.serde.Marshal(h.state)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %v", err)
	}

	_, err = h.context.SetState(map[string][]byte{
		KeyAddress: buff,
	})
	if err != nil {
		return fmt.Errorf("failed to set state: %v", err)
	}

	return nil
}

func loadState(ctx *processor.Context, serde serde.Serde) (state, error) {
	results, err := ctx.GetState([]string{KeyAddress})
	if err != nil {
		return state{}, fmt.Errorf("failed to get state: %v", err)
	}

	// In case the state is not initialized we use a new initialized state
	s := newState()

	result, found := results[KeyAddress]
	if found {
		err = serde.Unmarshal(result, &s)
		if err != nil {
			return state{}, fmt.Errorf("failed to unmarshal state: %v", err)
		}
	}

	return s, nil
}

func hexdigest(str string) string {
	hash := sha512.New()

	hash.Write([]byte(str))
	buff := hash.Sum(nil)

	return strings.ToLower(hex.EncodeToString(buff))
}
