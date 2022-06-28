package payload

import (
	"encoding/json"
	"fmt"
	"ledger/financial/state"
)

// FinancialPayload represents the payload sent in a financial transaction
type FinancialPayload struct {
	Action          string
	CreatePayload   CreatePayload
	TransferPayload TransferPayload
}

// CreatePayload defines the sub-payload "create"
type CreatePayload struct {
	WalletID state.WalletID
}

// TransferPayload defines the sub-payload "transfer"
type TransferPayload struct {
	FromWallet state.WalletID
	ToWallet   state.WalletID
	Amount     uint64
	Signature  string // hex encoded signature
}

// FromBytes unmarshals the buffer into a payload
func FromBytes(buff []byte) (FinancialPayload, error) {
	var payload FinancialPayload

	err := json.Unmarshal(buff, &payload)
	if err != nil {
		return FinancialPayload{}, fmt.Errorf("failed to unmarshal payload: %v", err)
	}

	return payload, nil
}

// ToBytes marshalls the payload into bytes
func (p FinancialPayload) ToBytes() ([]byte, error) {
	buff, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("failed to json marshal: %v", err)
	}

	return buff, nil
}
