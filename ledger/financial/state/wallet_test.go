package state

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWallet(t *testing.T) {
	wallets := Wallets{}

	w1 := Wallet{
		ID:      "aa",
		Balance: 10,
		State:   WalletOK,
	}

	w2 := Wallet{
		ID:      "bb",
		Balance: 20,
		State:   WalletBlacklisted,
	}

	wallets.Set(w1)
	require.Len(t, wallets, 1)

	wallets.Set(w1)
	require.Len(t, wallets, 1)

	wallets.Set(w2)
	require.Len(t, wallets, 2)

	wallets.Set(w2)
	require.Len(t, wallets, 2)

	wallets.Del(w1.ID)
	require.Len(t, wallets, 1)

	wallets.Del(w1.ID)
	require.Len(t, wallets, 1)

	_, found := wallets.Get(w1.ID)
	require.False(t, found)

	w, found := wallets.Get(w2.ID)
	require.True(t, found)
	require.Equal(t, w2, w)

	wallets.Del(w2.ID)
	require.Len(t, wallets, 0)

	_, found = wallets.Get(w2.ID)
	require.False(t, found)
}
