package controller

import (
	"go.dedis.ch/dela/cli"
	"go.dedis.ch/dela/cli/node"
)

// smcctl implements node.Initializer
type smcctl struct{}

// NewSmcController returns a new SMC initializer
func NewSmcController() node.Initializer {
	return smcctl{}
}

// SetCommands creates the CLI commands for the SMC part
func (s smcctl) SetCommands(builder node.Builder) {
	cmd := builder.SetCommand("smc")
	cmd.SetDescription("SMC service administration")

	sub := cmd.SetSubCommand("createkeys")
	sub.SetDescription("create key pair for reencryption")
	sub.SetAction(createKeyPairAction)

	sub = cmd.SetSubCommand("reveal")
	sub.SetDescription("reveal a reencrypted message")
	sub.SetFlags(
		cli.StringFlag{
			Name:  "xhatenc",
			Usage: "the reencrypted key as <hex(xhatenc)>",
		},
		cli.StringFlag{
			Name:  "dkgpub",
			Usage: "the DKG public key as <hex(dkgpub)>",
		},
		cli.StringFlag{
			Name:  "encrypted",
			Usage: "the encrypted string, as <hex(K)>:<hex(C1):<hex(C2):...>",
		},
		cli.StringFlag{
			Name:  "privk",
			Usage: "drop me if you can",
		},
	)
	sub.SetAction(revealAction)
}

// OnStart implements node.Initializer. It creates and registers a pedersen DKG.
func (s smcctl) OnStart(ctx cli.Flags, inj node.Injector) error {
	return nil
}

// OnStop implements node.Initializer.
func (s smcctl) OnStop(node.Injector) error {
	return nil
}
