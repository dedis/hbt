package state

type walletState int

// defines the states of a wallet
const (
	WalletOK = iota
	WalletBlacklisted
	WalletPending
)

// Wallet represents a main or sub-wallet
type Wallet struct {
	ID      WalletID
	Balance uint64
	State   walletState
}

// GetBalance returns the wallet's balance
func (w Wallet) GetBalance() uint64 {
	return w.Balance
}

// Wallets represents a list of wallets. A map would be more efficient but
// unfortunately not deterministic.
type Wallets []Wallet

// Get returns the wallet corresponding to the ID, or false.
func (ws Wallets) Get(id WalletID) (Wallet, bool) {
	for _, w := range ws {
		if w.ID == id {
			return w, true
		}
	}

	return Wallet{}, false
}

// Set sets a new wallets on the list. It ensures uniqueness of wallets based on
// their ID. This is to mimic a map. So, if a wallet already exists with the
// wallet ID, it will be replaced by the provided wallet.
func (ws *Wallets) Set(w1 Wallet) {
	for i, w2 := range *ws {
		if w2.ID == w1.ID {
			(*ws)[i] = w1
			return
		}
	}

	*ws = append(*ws, w1)
}

// Del remove a wallet based on its ID. It assumes that no two wallets have the
// same ID.
func (ws *Wallets) Del(id WalletID) bool {
	for i, w := range *ws {
		if w.ID == id {
			*ws = append((*ws)[:i], (*ws)[i+1:]...)
			return true
		}
	}

	return false
}
