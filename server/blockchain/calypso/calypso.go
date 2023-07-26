// Package calypso implements a Calypso contract that can advertise Secret
// Management Committees and deal with secrets.
//
// TODO: its information will be represented in the store as follows :
// SMC:C:{SMC public key} -> smc_roster (list of comma-separated host:port addresses)
// SMC:S:{secret name} -> secret value (encrypted with the SMC public key)
// SMC:R:{secret name} -> H(SMC public key, secret, client's public key)
package calypso

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net"
	"sort"
	"strings"

	"go.dedis.ch/dela"
	"go.dedis.ch/dela/core/access"
	"go.dedis.ch/dela/core/execution"
	"go.dedis.ch/dela/core/execution/native"
	"go.dedis.ch/dela/core/store"
	"golang.org/x/xerrors"
)

// commands defines the commands of the calypso contract.
// This interface helps in testing the contract.
type commands interface {
	advertiseSmc(snap store.Snapshot, step execution.Step) error
	deleteSmc(snap store.Snapshot, step execution.Step) error
	listSmc(snap store.Snapshot) error

	createSecret(snap store.Snapshot, step execution.Step) error
	listSecrets(snap store.Snapshot, step execution.Step) error
	revealSecret(snap store.Snapshot, step execution.Step) error
}

const (
	// ContractName is the name of the contract.
	ContractName = "go.dedis.ch/calypso.SMC"

	// KeyArg is the argument's name in the transaction that contains the
	// public key of the SMC to update.
	KeyArg = "calypso:smc_key"

	// RosterArg is the argument's name in the transaction that contains the
	// roster to associate with a given public key.
	RosterArg = "calypso:smc_roster"

	// SecretNameArg is the argument's name in the transaction that contains
	// the name of the secret to be published on the blockchain.
	SecretNameArg = "calypso:secret_name"

	// SecretArg is the argument's name in the transaction that contains the
	// secret to be published on the blockchain.
	SecretArg = "calypso:secret_value"

	// PubKeyArg is the argument's name in the transaction that contains the
	// public key to be used to re-encrypt the secret (and thus reveal it).
	PubKeyArg = "calypso:pub_key"

	// CmdArg is the argument's name to indicate the kind of command we want to
	// run on the contract. Should be one of the Command type.
	CmdArg = "calypso:command"

	// credentialAllCommand defines the credential command that is allowed to
	// perform all commands.
	credentialAllCommand = "all"
)

// Command defines a type of command for the value contract
type Command string

const (
	// CmdAdvertiseSmc defines the command to advertise a SMC
	CmdAdvertiseSmc Command = "ADVERTISE_SMC"

	// CmdDeleteSmc defines a command to delete a SMC
	CmdDeleteSmc Command = "DELETE_SMC"

	// CmdListSmc defines a command to list all SMCs (not deleted) so far.
	CmdListSmc Command = "LIST_SMC"

	// CmdCreateSecret defines a command to create a new secret.
	CmdCreateSecret Command = "CREATE_SECRET"

	// CmdListSecrets defines a command to list secrets for a SMC.
	CmdListSecrets Command = "LIST_SECRETS"

	// CmdRevealSecret defines a command to reveal a secret.
	CmdRevealSecret Command = "REVEAL_SECRET"
)

// Common error messages
const (
	notFoundInTxArg = "'%s' not found in tx arg"
)

// NewCreds creates new credentials for a value contract execution. We might
// want to use in the future a separate credential for each command.
func NewCreds(id []byte) access.Credential {
	return access.NewContractCreds(id, ContractName, credentialAllCommand)
}

// RegisterContract registers the value contract to the given execution service.
func RegisterContract(exec *native.Service, c Contract) {
	exec.Set(ContractName, c)
}

// Contract is a simple smart contract that allows one to handle the storage by
// performing create / list / update / delete operations on SMCs and secrets.
//
// - implements native.Contract
type Contract struct {
	// index contains all the keys set (and not deleteSmc) by this contract so far
	index map[string]struct{}

	// secrets contains a mapping between the keys and their associated secrets
	secrets map[string][]string

	// access is the access control service managing this smart contract
	access access.Service

	// accessKey is the access identifier allowed to use this smart contract
	accessKey []byte

	// cmd provides the commands that can be executed by this smart contract
	cmd commands

	// printer is the output used by the READ and LIST commands
	printer io.Writer
}

// NewContract creates a new Value contract
func NewContract(aKey []byte, srvc access.Service) Contract {
	contract := Contract{
		index:     map[string]struct{}{},
		secrets:   map[string][]string{},
		access:    srvc,
		accessKey: aKey,
		printer:   infoLog{},
	}

	contract.cmd = calypsoCommand{Contract: &contract}

	return contract
}

// Execute implements native.Contract. It runs the appropriate command.
func (c Contract) Execute(snap store.Snapshot, step execution.Step) error {
	creds := NewCreds(c.accessKey)

	err := c.access.Match(snap, creds, step.Current.GetIdentity())
	if err != nil {
		return xerrors.Errorf("identity not authorized: %v (%v)",
			step.Current.GetIdentity(), err)
	}

	cmd := step.Current.GetArg(CmdArg)
	if len(cmd) == 0 {
		return xerrors.Errorf(notFoundInTxArg, CmdArg)
	}

	switch Command(cmd) {
	case CmdAdvertiseSmc:
		err := c.cmd.advertiseSmc(snap, step)
		if err != nil {
			return xerrors.Errorf("failed to ADVERTISE_SMC: %v", err)
		}
	case CmdDeleteSmc:
		err := c.cmd.deleteSmc(snap, step)
		if err != nil {
			return xerrors.Errorf("failed to DELETE_SMC: %v", err)
		}
	case CmdListSmc:
		err := c.cmd.listSmc(snap)
		if err != nil {
			return xerrors.Errorf("failed to LIST_SMC: %v", err)
		}
	case CmdCreateSecret:
		err := c.cmd.createSecret(snap, step)
		if err != nil {
			return xerrors.Errorf("failed to CREATE_SECRET: %v", err)
		}
	case CmdListSecrets:
		err := c.cmd.listSecrets(snap, step)
		if err != nil {
			return xerrors.Errorf("failed to LIST_SECRETS: %v", err)
		}
	case CmdRevealSecret:
		err := c.cmd.revealSecret(snap, step)
		if err != nil {
			return xerrors.Errorf("failed to REVEAL_SECRET: %v", err)
		}
	default:
		return xerrors.Errorf("unknown command: %s", cmd)
	}

	return nil
}

// calypsoCommand implements the commands of the value contract
//
// - implements commands
type calypsoCommand struct {
	*Contract
}

// advertiseSmc implements commands. It performs the ADVERTISE_SMC command.
// It can advertise a new SMC or update the roster of an existing one.
func (c calypsoCommand) advertiseSmc(snap store.Snapshot, step execution.Step) error {
	key := step.Current.GetArg(KeyArg)
	if len(key) == 0 {
		return xerrors.Errorf(notFoundInTxArg, KeyArg)
	}

	roster := step.Current.GetArg(RosterArg)
	if len(roster) == 0 {
		return xerrors.Errorf(notFoundInTxArg, RosterArg)
	}

	nodeList := strings.Split(string(roster), ",")
	for _, r := range nodeList {
		_, _, err := net.SplitHostPort(r)
		if err != nil {
			return xerrors.Errorf("invalid node '%s' in roster: %v", r, err)
		}
	}

	currentRoster, err := snap.Get(key)
	if err == nil && len(currentRoster) > 0 {
		// the SMC already exists, we need to verify the intersection between
		// the new roster and the old one. There must be at least a threshold
		// of nodes in the intersection

		e := validateRosterUpdate(currentRoster, roster)
		if e != nil {
			return xerrors.Errorf("roster validation failed: %v", e)
		}
	}

	err = snap.Set(key, roster) // DKG public key => roster
	if err != nil {
		return xerrors.Errorf("failed to set roster: %v", err)
	}

	c.index[string(key)] = struct{}{}
	c.secrets[string(key)] = []string{}

	dela.Logger.Info().Str("contract", ContractName).Msgf("setting %x=%s", key, roster)

	return nil
}

// validateRosterUpdate verifies that the new roster has sufficient overlap
// with the old roster. It returns an error if the new roster is not valid.
func validateRosterUpdate(oldRoster []byte, newRoster []byte) error {
	oldRosterList := strings.Split(string(oldRoster), ",")
	newRosterList := strings.Split(string(newRoster), ",")

	sort.Strings(oldRosterList)
	sort.Strings(newRosterList)

	// verify there's overlap between old roster and new roster
	overlap := intersectSortedRosters(oldRosterList, newRosterList)
	thr := len(oldRosterList) - (len(oldRosterList)-1)/3
	if overlap < thr {
		return xerrors.Errorf(
			"new roster does not overlap enough with current roster (%d < %d)",
			overlap, thr)
	}

	return nil
}

// intersectSortedRosters returns the # of elements in common between 2 rosters
// Its behaviour is undefined if the rosters are not sorted.
func intersectSortedRosters(oldRoster []string, newRoster []string) int {
	overlap := 0

	oldIdx := 0
	newIdx := 0
	for oldIdx < len(oldRoster) && newIdx < len(newRoster) {
		o, n := oldRoster[oldIdx], newRoster[newIdx]

		switch {
		case o < n:
			oldIdx++

		case o > n:
			newIdx++

		case o == n:
			oldIdx++
			newIdx++
			overlap++
		}
	}

	return overlap
}

// deleteSmc implements commands. It performs the DELETE_SMC command
func (c calypsoCommand) deleteSmc(snap store.Snapshot, step execution.Step) error {
	key := step.Current.GetArg(KeyArg)
	if len(key) == 0 {
		return xerrors.Errorf(notFoundInTxArg, KeyArg)
	}

	err := snap.Delete(key)
	if err != nil {
		return xerrors.Errorf("failed to deleteSmc key '%x': %v", key, err)
	}

	// DKG => roster
	// DKG => [secret1_key, secret2, secret3, ...]
	// secret1_key => secret1_encrypted_value
	for _, secret := range c.secrets[string(key)] {
		dela.Logger.Info().
			Msgf("Deleting secret '%s' that depended on deleted SMC '%s'", secret, key)

		err = snap.Delete([]byte(secret))
		if err != nil {
			dela.Logger.Warn().
				Msgf("Could not delete secret '%s', "+
					"orphaned by deleted SMC '%s'", secret, key)
		}
	}

	delete(c.index, string(key))
	delete(c.secrets, string(key))

	return nil
}

// listSmc implements commands. It performs the LIST_SMC command
func (c calypsoCommand) listSmc(snap store.Snapshot) error {
	res := []string{}

	for k := range c.index {
		v, err := snap.Get([]byte(k))
		if err != nil {
			return xerrors.Errorf("failed to get key '%s': %v", k, err)
		}

		res = append(res, fmt.Sprintf("%x=%s", k, v))
	}

	sort.Strings(res)
	fmt.Fprint(c.printer, strings.Join(res, ","))

	return nil
}

// createSecret implements commands. It performs the CREATE_SECRET command
func (c calypsoCommand) createSecret(snap store.Snapshot, step execution.Step) error {

	key := step.Current.GetArg(KeyArg)
	if len(key) == 0 {
		return xerrors.Errorf(notFoundInTxArg, KeyArg)
	}

	name := step.Current.GetArg(SecretNameArg)
	if len(name) == 0 {
		return xerrors.Errorf(notFoundInTxArg, SecretNameArg)
	}

	secret := step.Current.GetArg(SecretArg)
	if len(secret) == 0 {
		return xerrors.Errorf(notFoundInTxArg, SecretArg)
	}

	_, ok := c.index[string(key)]
	if !ok {
		return xerrors.Errorf("'%s' was not found among the SMCs", key)
	}

	err := snap.Set(name, secret)
	if err != nil {
		return xerrors.Errorf("failed to set secret: %v", err)
	}

	c.secrets[string(key)] = append(c.secrets[string(key)], string(name))

	dela.Logger.Info().
		Str("contract", ContractName).
		Msgf("setting secret %x=%s", name, secret)

	return nil
}

// listSecrets implements commands. It performs the LIST_SECRETS command
func (c calypsoCommand) listSecrets(snap store.Snapshot, step execution.Step) error {
	res := []string{}

	key := step.Current.GetArg(KeyArg)
	if len(key) == 0 {
		return xerrors.Errorf(notFoundInTxArg, KeyArg)
	}

	_, found := c.secrets[string(key)]
	if !found {
		return xerrors.Errorf("SMC not found: %s", key)
	}

	for _, k := range c.secrets[string(key)] {
		v, err := snap.Get([]byte(k))
		if err != nil {
			return xerrors.Errorf("failed to get key '%s': %v", k, err)
		}

		res = append(res, fmt.Sprintf("%s=%s", k, v))
	}

	sort.Strings(res)
	fmt.Fprint(c.printer, strings.Join(res, ","))

	return nil
}

// revealSecret implements commands. It performs the REVEAL_SECRET command
func (c calypsoCommand) revealSecret(snap store.Snapshot, step execution.Step) error {

	smcKey := step.Current.GetArg(KeyArg)
	if len(smcKey) == 0 {
		return xerrors.Errorf(notFoundInTxArg, KeyArg)
	}

	name := step.Current.GetArg(SecretNameArg)
	if len(name) == 0 {
		return xerrors.Errorf(notFoundInTxArg, SecretNameArg)
	}

	clientPubKey := step.Current.GetArg(PubKeyArg)
	if len(clientPubKey) == 0 {
		return xerrors.Errorf(notFoundInTxArg, PubKeyArg)
	}

	smcSecrets, ok := c.secrets[string(smcKey)]
	if !ok {
		return xerrors.Errorf("'%s' was not found among the SMCs", smcKey)
	}

	found := false
	for _, s := range smcSecrets {
		if s == string(name) {
			found = true
			break
		}
	}

	if !found {
		return xerrors.Errorf(
			"'%s' was not found among the secrets of the smc (%v)",
			name, smcKey)
	}

	secret, err := snap.Get(name)
	if err != nil {
		return xerrors.Errorf("failed to get secret '%s': %v", name, err)
	}

	// Generate authorization hash
	h := sha256.New()
	h.Write(smcKey)
	h.Write(secret)
	h.Write(clientPubKey)
	hash := h.Sum(nil)

	// TODO: this should be done according to the store specification at the
	// top of this file, but the prefixing logic has not yet been implemented.
	err = snap.Set(hash, hash)
	if err != nil {
		return xerrors.Errorf("failed to persist secret reveal: %v", err)
	}

	dela.Logger.Info().
		Str("contract", ContractName).
		Msgf("revealed secret %x to %s", name, step.Current.GetIdentity())

	return nil
}

// infoLog defines an output using zerolog
//
// - implements io.writer
type infoLog struct{}

func (h infoLog) Write(p []byte) (int, error) {
	dela.Logger.Info().Msg(string(p))

	return len(p), nil
}
