package main

import (
	"fmt"
	"ledger/financial/serde"
	"net/url"
	"os"

	"github.com/hyperledger/sawtooth-sdk-go/logging"
	"github.com/jessevdk/go-flags"
	flag "github.com/jessevdk/go-flags"
)

var defaultSerde = serde.JSON{}

const (
	// BatchSubmitAPI represents the REST endpoint to send batches
	BatchSubmitAPI string = "batches"

	// ContentTypeOctetStream represents the octet stream content type
	ContentTypeOctetStream string = "application/octet-stream"
)

// Opts defines the main CLI options
type Opts struct {
	Verbose        []bool `short:"v" long:"verbose" description:"Increase verbosity"`
	PrivateKeyFile string `short:"k" long:"privatekeyfile" description:"Filepath of the private key that is allowed sign batches" required:"true"`
}

type command interface {
	flag.Commander

	register(*flag.Command) error
}

var logger *logging.Logger = logging.Get()

func main() {
	var opts Opts
	parser := flags.NewParser(&opts, flags.Default)
	parser.Command.Name = "financial"

	url := url.URL{
		Scheme: "http",
		Host:   "localhost:8008",
	}

	// needed to already have the private key file
	parser.Parse()

	batcher, err := newBatcher(url, opts.PrivateKeyFile)
	if err != nil {
		logger.Errorf("failed to create batcher: %v", err)
		os.Exit(1)
	}

	commands := []command{
		newCmdAdd(batcher),
	}

	for _, command := range commands {
		err := command.register(parser.Command)
		if err != nil {
			logger.Errorf("failed to register command: %v", err)
			os.Exit(1)
		}
	}

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
		os.Exit(1)
	}

	switch len(opts.Verbose) {
	case 2:
		logger.SetLevel(logging.DEBUG)
	case 1:
		logger.SetLevel(logging.INFO)
	default:
		logger.SetLevel(logging.WARN)
	}

	if parser.Command.Active == nil {
		fmt.Println("no command provided")
		os.Exit(2)
	}
}
