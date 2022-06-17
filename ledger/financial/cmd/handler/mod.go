package main

import (
	"ledger/financial/handler"
	"ledger/financial/serde"
	"os"
	"syscall"

	"github.com/hyperledger/sawtooth-sdk-go/logging"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	flags "github.com/jessevdk/go-flags"
)

var defaultSerde = serde.JSON{}

// Opts defines the CLI options
type Opts struct {
	Verbose []bool `short:"v" long:"verbose" description:"Increase verbosity"`
	Connect string `short:"C" long:"connect" description:"Validator component endpoint to connect to" default:"tcp://localhost:4004"`
}

func main() {
	var opts Opts

	logger := logging.Get()

	parser := flags.NewParser(&opts, flags.Default)

	remaining, err := parser.Parse()
	if err != nil {
		flagsErr, ok := err.(*flags.Error)
		if ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}
		logger.Errorf("failed to parse args: %v", err)
		os.Exit(2)
	}

	if len(remaining) != 0 {
		logger.Errorf("unrecognized arguments: %v", remaining)
	}

	endpoint := opts.Connect

	switch len(opts.Verbose) {
	case 2:
		logger.SetLevel(logging.DEBUG)
	case 1:
		logger.SetLevel(logging.INFO)
	default:
		logger.SetLevel(logging.WARN)
	}

	logger.Debugf("command line arguments: %v\n", os.Args)
	logger.Debugf("verbose = %v\n", len(opts.Verbose))
	logger.Debugf("endpoint = %v\n", endpoint)

	processor := processor.NewTransactionProcessor(endpoint)
	processor.AddHandler(handler.NewFinancialHandler(defaultSerde))
	processor.ShutdownOnSignal(syscall.SIGINT, syscall.SIGTERM)

	err = processor.Start()
	if err != nil {
		logger.Errorf("failed to start: %v", err)
		os.Exit(1)
	}
}
