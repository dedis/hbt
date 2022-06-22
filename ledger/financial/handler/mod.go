package handler

import (
	"fmt"
	"ledger/financial/payload"
	"ledger/financial/serde"
	"ledger/financial/state"

	"github.com/hyperledger/sawtooth-sdk-go/logging"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/processor_pb2"
)

// defines the transaction family constants
const (
	FamilyName    = "hbt_financial"
	FamilyVersion = "1.0"
)

var logger *logging.Logger = logging.Get()

// NewFinancialHandler returns a new initialized financial handler
func NewFinancialHandler(serde serde.Serde) processor.TransactionHandler {
	return &FinancialHandler{
		serde: serde,
	}
}

// FinancialHandler defines the handler of the financial smart contract
//
// - implement processor.TransactionHandler
type FinancialHandler struct {
	serde serde.Serde
}

// FamilyName implement processor.TransactionHandler. It should return the name
// of the transaction family that this handler can process, e.g. "intkey"
func (f *FinancialHandler) FamilyName() string {
	return FamilyName
}

// FamilyVersions implement processor.TransactionHandler. It should return the
// versions of the transaction family that this handler can process. Eg.,
// ["1.0", "1.5"]
func (f *FinancialHandler) FamilyVersions() []string {
	return []string{FamilyVersion}
}

// Namespaces implement processor.TransactionHandler. It should return a slice
// containing all the handler's namespaces, e.g. []string{"abcdef"}
func (f *FinancialHandler) Namespaces() []string {
	return []string{state.KeyAddress[:6]}
}

// Apply implement processor.TransactionHandler. It is the single method where
// all the business logic for a transaction family is defined. The method will
// be called by the transaction processor upon receiving a TpProcessRequest that
// the handler understands and will pass in the TpProcessRequest and an
// initialized instance of the Context type.
func (f *FinancialHandler) Apply(request *processor_pb2.TpProcessRequest,
	context *processor.Context) error {

	payload, err := payload.FromBytes(request.GetPayload())
	if err != nil {
		return fmt.Errorf("failed to get payload: %v", err)
	}

	logger.Debugf("trading tx %v: action: %s, walletID: %s", request.Signature,
		payload.Action, payload.WalletID)

	stateHandler := state.NewHandler(context, f.serde)

	switch payload.Action {
	case "create":
		err := stateHandler.CreateWallet(payload.WalletID)
		if err != nil {
			return &processor.InvalidTransactionError{
				Msg: fmt.Sprintf("failed to create wallet: %v", err),
			}
		}

		logger.Infof("wallet with id %q created", payload.WalletID)
	case "get":
		wallet, err := stateHandler.GetWallet(payload.WalletID)
		if err != nil {
			return &processor.InvalidTransactionError{
				Msg: fmt.Sprintf("failed to get wallet: %v", err),
			}
		}

		logger.Infof("wallet [%s] balance=%d", payload.WalletID, wallet.GetBalance())
	default:
		return &processor.InvalidTransactionError{
			Msg: fmt.Sprintf("invalid action: %q", payload.Action),
		}
	}

	return nil
}
