package payload

import (
	"encoding/json"
	"fmt"
)

// FinancialPayload represents the payload sent in a financial transaction
type FinancialPayload struct {
	Action   string
	WalletID string
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
