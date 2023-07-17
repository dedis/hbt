package controller

import (
	"go.dedis.ch/dela/cli"
	"go.dedis.ch/dela/cli/node"
	"go.dedis.ch/kyber/v3/suites"
	"go.dedis.ch/kyber/v3/util/key"
)

// smcctl is an initializer with a set of commands for the SMC. It only
// creates and injects a new SMC instance.
//
// - implements node.Initializer
type smcctl struct {
	kp *key.Pair
}

// NewSmcController returns a new SMC initializer
func NewSmcController() node.Initializer {
	return smcctl{
		kp: key.NewKeyPair(suites.MustFind("Ed25519")),
	}
}

// Build implements node.Initializer. In this case we don't need any command.
func (s smcctl) SetCommands(builder node.Builder) {
	cmd := builder.SetCommand("smc")
	cmd.SetDescription("SMC service administration")

	sub := cmd.SetSubCommand("createkeys")
	sub.SetDescription("create key pair for reencryption")
	sub.SetAction(builder.MakeAction(createKpAction{}))

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
	)
	sub.SetAction(builder.MakeAction(revealAction{}))
}

// OnStart implements node.Initializer. It creates and registers a pedersen DKG.
func (s smcctl) OnStart(ctx cli.Flags, inj node.Injector) error {
	return nil
}

// OnStop implements node.Initializer.
func (s smcctl) OnStop(node.Injector) error {
	return nil
}
