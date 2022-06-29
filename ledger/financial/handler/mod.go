package handler

import (
	"crypto/sha512"
	"fmt"
	"ledger/financial/payload"
	"ledger/financial/serde"
	"ledger/financial/state"
	"strconv"

	"github.com/hyperledger/sawtooth-sdk-go/logging"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/processor_pb2"
	"go.dedis.ch/kyber/v3/sign/schnorr"
	"go.dedis.ch/kyber/v3/suites"
)

// defines the transaction family constants
const (
	FamilyName    = "hbt_financial"
	FamilyVersion = "1.0"
)

var suite = suites.MustFind("ed25519")

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

	logger.Debugf("financial tx %v: action: %s", request.Signature, payload.Action)

	stateHandler, err := state.NewHandler(context, f.serde)
	if err != nil {
		return fmt.Errorf("faield to create state handler: %v", err)
	}

	switch payload.Action {
	case "create":
		createPayload := payload.CreatePayload

		err := stateHandler.CreateWallet(createPayload.WalletID)
		if err != nil {
			return &processor.InvalidTransactionError{
				Msg: fmt.Sprintf("failed to create wallet: %v", err),
			}
		}

		logger.Infof("wallet with id %q created", createPayload.WalletID)
	case "transfer":
		err := f.transfer(stateHandler, payload.TransferPayload)
		if err != nil {
			return &processor.InvalidTransactionError{
				Msg: fmt.Sprintf("failed to call transfer: %v", err),
			}
		}
	default:
		return &processor.InvalidTransactionError{
			Msg: fmt.Sprintf("invalid action: %q", payload.Action),
		}
	}

	return nil
}

func (f *FinancialHandler) transfer(state state.Handler,
	payload payload.TransferPayload) error {

	fromWallet, err := state.GetWallet(payload.FromWallet)
	if err != nil {
		return fmt.Errorf("failed to get 'from' wallet %q: %v", payload.FromWallet, err)
	}

	fromKey, err := payload.FromWallet.ToPoint()
	if err != nil {
		return fmt.Errorf("failed to key fromKey: %v", err)
	}

	hash := sha512.New()
	hash.Write([]byte(strconv.Itoa(int(payload.Amount))))
	hash.Write([]byte(payload.FromWallet))
	hash.Write([]byte(payload.ToWallet))

	err = schnorr.Verify(suite, fromKey, hash.Sum(nil), []byte(payload.Signature))
	if err != nil {
		return fmt.Errorf("signature verification failed: %v", err)
	}

	if fromWallet.GetBalance() < payload.Amount {
		return fmt.Errorf("insufficient funds: %d < %d", fromWallet.GetBalance(),
			payload.Amount)
	}

	toWallet, err := state.GetWallet(payload.ToWallet)
	if err != nil {
		return fmt.Errorf("failed to get 'from' wallet %q: %v", payload.ToWallet, err)
	}

	fromWallet.Balance -= payload.Amount
	toWallet.Balance += payload.Amount

	err = state.StoreState()
	if err != nil {
		return fmt.Errorf("failed to store state: %v", err)
	}

	return nil
}
