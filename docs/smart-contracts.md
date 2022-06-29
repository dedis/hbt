# Smart contracts

This is a working documentation about the smart contracts used in HBT.

## Financial

### Purpose

The financial smart contract handles transfer of tokens between two accounts.

### State

The financial smart contract holds a list of wallets, which are represented as
type `WalletID`. A `WalletID` is a hex-encoded public key, which is a point on
the ED25519 elliptic curve.

### Operation: Create

The create operation creates a new wallet, with an initial balance of 0. The
payload of this operation is the following:

```go
type CreatePayload struct {
	WalletID state.WalletID
}
```

No verification is performed.

### Operation: Transfer

The transfer operation transfers tokens between one wallet to another. The
payload of this operation is the following:

```go
type TransferPayload struct {
	FromWallet state.WalletID
	ToWallet   state.WalletID
	Amount     uint64
	Signature  string // hex encoded signature
}
```

The verification is as follow:

- Signature must be done by `FromWallet` on `SHA512(Amount || FromWallet ||
ToWallet)`.
- The balance of the FromWallet must be >= to Amount