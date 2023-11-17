package controller

import (
	"go.dedis.ch/hbt/server/blockchain/calypso"

	"go.dedis.ch/dela/cli"
	"go.dedis.ch/dela/cli/node"
	"go.dedis.ch/dela/core/access"
	"go.dedis.ch/dela/core/execution/native"
	"golang.org/x/xerrors"
)

// miniController is a CLI initializer to register the value contract
//
// - implements node.Initializer
type miniController struct {
}

// NewController creates a new minimal controller for the value contract.
func NewController() node.Initializer {
	return miniController{}
}

// SetCommands implements node.Initializer.
func (miniController) SetCommands(_ node.Builder) {
	// For now, no need to do anything here.
}

// OnStart implements node.Initializer. It registers the value contract.
func (m miniController) OnStart(_ cli.Flags, inj node.Injector) error {
	var access access.Service
	err := inj.Resolve(&access)
	if err != nil {
		return xerrors.Errorf("failed to resolve access service: %v", err)
	}

	var exec *native.Service
	err = inj.Resolve(&exec)
	if err != nil {
		return xerrors.Errorf("failed to resolve native service: %v", err)
	}

	contract := calypso.NewContract(access)

	calypso.RegisterContract(exec, contract)

	return nil
}

// OnStop implements node.Initializer.
func (miniController) OnStop(_ node.Injector) error {
	return nil
}
