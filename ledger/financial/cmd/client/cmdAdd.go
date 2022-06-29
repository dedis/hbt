package main

import (
	"fmt"
	"ledger/financial/payload"
	"ledger/financial/state"

	flag "github.com/jessevdk/go-flags"
)

func newCmdAdd(batcher batcher) command {
	return &cmdAdd{
		batcher: batcher,
	}
}

// cmdAdd defines a command to add a new wallet
//
// - implements command
type cmdAdd struct {
	Test        string `long:"test" description:"This is a test arg"`
	KeyFilePath string `short:"k" long:"keyfilepath" description:"Filepath of the private key" default:""`
	ID          string `long:"id" description:"the wallet ID" required:"true"`

	batcher batcher
}

// register implements command
func (a *cmdAdd) register(command *flag.Command) error {
	_, err := command.AddCommand("add", "adds a wallet", "adds a waller", a)
	if err != nil {
		return fmt.Errorf("faield to add command: %v", err)
	}

	return nil
}

// Execute implements flag.Commander
func (a cmdAdd) Execute(args []string) error {

	payload := payload.FinancialPayload{
		Action: "create",
		CreatePayload: payload.CreatePayload{
			WalletID: state.WalletID(a.ID),
		},
	}

	res, err := a.batcher.sendTransaction(payload)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %v", err)
	}

	fmt.Println("res:", res)

	return nil
}
