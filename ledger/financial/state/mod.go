package state

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"strings"

	"ledger/financial/serde"

	"github.com/hyperledger/sawtooth-sdk-go/processor"
)

// KeyAddress defines the key address of the financial smart contract state
var KeyAddress = hexdigest("hbt_financial")[:70]

// newState returns a new initialized state
func newState() state {
	return state{
		wallets: make(map[string]Wallet),
	}
}

// state defines the data stored by the smart contract
type state struct {
	wallets map[string]Wallet
}

// Wallet represents a main or sub-wallet
type Wallet struct {
	balance uint64
}

// GetBalance returns the wallet's balance
func (w Wallet) GetBalance() uint64 {
	return w.balance
}

// NewHandler returns a new initialized Handler
func NewHandler(ctx *processor.Context, serde serde.Serde) Handler {
	return Handler{
		context: ctx,
		serde:   serde,
	}
}

// Handler handles addressing and (de)serialization of the state
type Handler struct {
	context *processor.Context
	serde   serde.Serde
}

// CreateWallet creates a new wallet with walletID. Returns an error if the
// wallet already exists.
func (h Handler) CreateWallet(walletID string) error {
	state, err := h.loadState()
	if err != nil {
		return fmt.Errorf("failed to load state: %v", err)
	}

	_, found := state.wallets[walletID]
	if found {
		return fmt.Errorf("wallet with id %q already exists", walletID)
	}

	state.wallets[walletID] = Wallet{
		balance: 0,
	}

	err = h.storeState(state)
	if err != nil {
		return fmt.Errorf("failed to store state: %v", err)
	}

	return nil
}

// GetWallet returns the wallet corresponding to the walletID, or an error if
// not found.
func (h Handler) GetWallet(walletID string) (Wallet, error) {
	state, err := h.loadState()
	if err != nil {
		return Wallet{}, fmt.Errorf("failed to load state: %v", err)
	}

	wallet, found := state.wallets[walletID]
	if !found {
		return Wallet{}, fmt.Errorf("wallet with id %q already exists", walletID)
	}

	return wallet, nil
}

func (h Handler) loadState() (state, error) {
	results, err := h.context.GetState([]string{KeyAddress})
	if err != nil {
		return state{}, fmt.Errorf("failed to get state: %v", err)
	}

	// In case the state is not initialized we use a new initialized state
	s := newState()

	result, found := results[KeyAddress]
	if found {
		err = h.serde.Unmarshal(result, &s)
		if err != nil {
			return state{}, fmt.Errorf("failed to unmarshal state: %v", err)
		}
	}

	return s, nil
}

func (h Handler) storeState(state state) error {
	buff, err := h.serde.Marshal(state)
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

func hexdigest(str string) string {
	hash := sha512.New()

	hash.Write([]byte(str))
	buff := hash.Sum(nil)

	return strings.ToLower(hex.EncodeToString(buff))
}
