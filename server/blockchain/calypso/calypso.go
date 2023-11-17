// Package calypso implements a Calypso contract that can advertise Secret
// Management Committees, deal with secrets and audit their access.
//
// Its information will be represented in the KV store as follows :
// CALYR {SMC pub key} -> {pub key, smc_roster (list of comma-separated host:port addresses)}
// CALYS {secret name} -> secret value (encrypted with the SMC public key)
// CALYL {secret name} -> H(SMC public key, secret, client's public key)
// CALYA {H(SMC public key, secret, client's public key)} -> client's public key
package calypso

import (
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
	"go.dedis.ch/dela/core/store/prefixed"
	"go.dedis.ch/dela/crypto"
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

	listAuditLogs(snap store.Snapshot, step execution.Step) error
}

const (
	// ContractUID is the unique (4-bytes) identifier of the contract, it is
	// used to prefix keys in the K/V store and by DARCs for access control.
	ContractUID = "CALY"

	// ContractName is the name of the contract.
	ContractName = "go.dedis.ch/calypso.SMC"

	// SmcPublicKeyArg is the argument's name in the transaction that contains
	// the public key of the SMC to update.
	SmcPublicKeyArg = "calypso:smc_key"

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

	// PrefixSmcRosterKeys prefixed store keys contain the roster of the SMC.
	PrefixSmcRosterKeys = ContractUID + "R"

	// PrefixSecretKeys prefixed store keys contain the secret.
	PrefixSecretKeys = ContractUID + "S"

	// PrefixListKeys prefixed store keys contain the list of audit keys
	// that had access to the secret.
	// e.g. [SMCL|Secret] => [SMCA1, SMCA2, SMCA3, ...]
	PrefixListKeys = ContractUID + "L"

	// PrefixAccessKeys prefixed store keys contain the public key of the secret reader
	// e.g. [SMCA|Hash(...)] => PubKey
	PrefixAccessKeys = ContractUID + "A"

	// errorKeyNotFoundInSmcs is used in error messages of this module
	errorKeyNotFoundInSmcs = "'%s' was not found among the SMCs"

	// the length of the access token (depending on the sha used)
	accessTokenNbBytes = 32
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

	// CmdListAuditLog defines a command to list audit logs.
	CmdListAuditLog Command = "LIST_AUDIT_LOG"
)

// Common error messages
const (
	notFoundInTxArg = "'%s' not found in tx arg"
)

// smcPubKey contains the SMC public key.
type smcPubKey string

type secretSet map[string]struct{}

func (s secretSet) addSecret(name string) {
	s[name] = struct{}{}
}

// NewCreds creates new credentials for a value contract execution. We might
// want to use in the future a separate credential for each command.
func NewCreds() access.Credential {
	return access.NewContractCreds([]byte(ContractUID), ContractName, credentialAllCommand)
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
	// secrets contains a mapping between the SMC and their associated secrets
	// k=SMC pub key, v=set of secret names
	secrets map[smcPubKey]secretSet

	// access is the access control service managing this smart contract
	access access.Service

	// cmd provides the commands that can be executed by this smart contract
	cmd commands

	// printer is the output used by the READ and LIST commands
	printer io.Writer
}

// NewContract creates a new Calypso contract
func NewContract(srvc access.Service) Contract {
	contract := Contract{
		secrets: map[smcPubKey]secretSet{},
		access:  srvc,
		printer: infoLog{},
	}

	contract.cmd = calypsoCommand{Contract: &contract}

	return contract
}

// Execute implements native.Contract. It checks that command is formed correctly before running it.
func (c Contract) Execute(snap store.Snapshot, step execution.Step) error {
	creds := NewCreds()

	err := c.access.Match(snap, creds, step.Current.GetIdentity())
	if err != nil {
		return xerrors.Errorf("identity not authorized: %v (%v)",
			step.Current.GetIdentity(), err)
	}

	cmd := step.Current.GetArg(CmdArg)
	if len(cmd) == 0 {
		return xerrors.Errorf(notFoundInTxArg, CmdArg)
	}

	return c.ExecuteCommand(snap, step, cmd)
}

// ExecuteCommand executes the appropriate command.
func (c Contract) ExecuteCommand(snap store.Snapshot, step execution.Step, cmd []byte) error {
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
	case CmdListAuditLog:
		err := c.cmd.listAuditLogs(snap, step)
		if err != nil {
			return xerrors.Errorf("failed to LIST_AUDIT_LOG: %v", err)
		}
	default:
		return xerrors.Errorf("unknown command: %s", cmd)
	}

	return nil
}

// UID returns the unique 4-bytes contract identifier.
//
// - implements native.Contract
func (c Contract) UID() string {
	return ContractUID
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
	key := step.Current.GetArg(SmcPublicKeyArg)
	if len(key) == 0 {
		return xerrors.Errorf(notFoundInTxArg, SmcPublicKeyArg)
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

	currentRoster, err := getSmcRoster(snap, key)
	if err == nil && len(currentRoster) > 0 {
		// the SMC already exists, we need to verify the intersection between
		// the new roster and the old one. There must be at least a threshold
		// of nodes in the intersection

		e := validateRosterUpdate(currentRoster, roster)
		if e != nil {
			return xerrors.Errorf("roster validation failed: %v", e)
		}
	}

	err = setSmcRoster(snap, key, roster) // DKG public key => roster
	if err != nil {
		return xerrors.Errorf("failed to set roster: %v", err)
	}

	c.secrets[smcPubKey(key)] = secretSet{}

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
	key := step.Current.GetArg(SmcPublicKeyArg)
	if len(key) == 0 {
		return xerrors.Errorf(notFoundInTxArg, SmcPublicKeyArg)
	}

	k := prefixed.NewPrefixedKey([]byte(PrefixSmcRosterKeys), key)

	err := snap.Delete(k)
	if err != nil {
		return xerrors.Errorf("failed to delete SMC with public key '%x': %v", key, err)
	}

	// DKG => roster
	// DKG => [secret1_key, secret2, secret3, ...]
	// secret1_key => secret1_encrypted_value
	for secret := range c.secrets[smcPubKey(key)] {
		dela.Logger.Info().
			Msgf("Deleting secret '%s' that depended on deleted SMC '%s'", secret, key)

		err = deleteSecret(snap, []byte(secret))
		if err != nil {
			dela.Logger.Warn().
				Msgf("Could not delete secret '%s', "+
					"orphaned by deleted SMC '%s'", secret, key)
		}
	}

	delete(c.secrets, smcPubKey(key))

	return nil
}

// listSmc implements commands. It performs the LIST_SMC command
func (c calypsoCommand) listSmc(snap store.Snapshot) error {
	res := []string{}

	for k := range c.secrets {
		v, err := getSmcRoster(snap, []byte(k))
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
	smcKey := step.Current.GetArg(SmcPublicKeyArg)
	if len(smcKey) == 0 {
		return xerrors.Errorf(notFoundInTxArg, SmcPublicKeyArg)
	}

	name := step.Current.GetArg(SecretNameArg)
	if len(name) == 0 {
		return xerrors.Errorf(notFoundInTxArg, SecretNameArg)
	}

	secret := step.Current.GetArg(SecretArg)
	if len(secret) == 0 {
		return xerrors.Errorf(notFoundInTxArg, SecretArg)
	}

	_, ok := c.secrets[smcPubKey(smcKey)]
	if !ok {
		return xerrors.Errorf(errorKeyNotFoundInSmcs, smcKey)
	}

	v, _ := getSecret(snap, name)
	if v != nil {
		return xerrors.Errorf("a secret named '%s' already exists", name)
	}

	err := setSecret(snap, name, secret)
	if err != nil {
		return xerrors.Errorf("failed to set secret: %v", err)
	}

	c.secrets[smcPubKey(smcKey)].addSecret(string(name))

	return nil
}

// listSecrets implements commands. It performs the LIST_SECRETS command
func (c calypsoCommand) listSecrets(snap store.Snapshot, step execution.Step) error {
	res := []string{}

	key := step.Current.GetArg(SmcPublicKeyArg)
	if len(key) == 0 {
		return xerrors.Errorf(notFoundInTxArg, SmcPublicKeyArg)
	}

	_, found := c.secrets[smcPubKey(key)]
	if !found {
		return xerrors.Errorf("SMC not found: %s", key)
	}

	for k := range c.secrets[smcPubKey(key)] {
		v, err := getSecret(snap, []byte(k))
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

	smcKey := step.Current.GetArg(SmcPublicKeyArg)
	if len(smcKey) == 0 {
		return xerrors.Errorf(notFoundInTxArg, SmcPublicKeyArg)
	}

	name := step.Current.GetArg(SecretNameArg)
	if len(name) == 0 {
		return xerrors.Errorf(notFoundInTxArg, SecretNameArg)
	}

	clientPubKey := step.Current.GetArg(PubKeyArg)
	if len(clientPubKey) == 0 {
		return xerrors.Errorf(notFoundInTxArg, PubKeyArg)
	}

	smcSecrets, ok := c.secrets[smcPubKey(smcKey)]
	if !ok {
		return xerrors.Errorf(errorKeyNotFoundInSmcs, smcKey)
	}

	_, found := smcSecrets[string(name)]

	if !found {
		return xerrors.Errorf(
			"'%s' was not found among the secrets of the smc (%v)",
			name, string(smcKey))
	}

	secret, err := getSecret(snap, name)
	if err != nil {
		return xerrors.Errorf("failed to get secret '%s': %v", name, err)
	}

	accessToken := computeAccessToken(smcKey, secret, clientPubKey)

	if hasSecretAccess(snap, accessToken) {
		return nil
	}

	err = setSecretAccess(snap, accessToken, clientPubKey)
	if err != nil {
		return xerrors.Errorf("failed to persist secret reveal: %v", err)
	}

	err = insertAuditLog(snap, name, accessToken)
	if err != nil {
		return xerrors.Errorf("failed to persist secret audit log: %v", err)
	}

	dela.Logger.Info().
		Str("contract", ContractName).
		Msgf("revealed secret %x to %s", name, step.Current.GetIdentity())

	return nil
}

// listAuditLogs implements commands. It performs the LIST_AUDIT_LOGS command
func (c calypsoCommand) listAuditLogs(snap store.Snapshot, step execution.Step) error {
	smcKey := step.Current.GetArg(SmcPublicKeyArg)
	if len(smcKey) == 0 {
		return xerrors.Errorf(notFoundInTxArg, SmcPublicKeyArg)
	}

	name := step.Current.GetArg(SecretNameArg)
	if len(name) == 0 {
		return xerrors.Errorf(notFoundInTxArg, SecretNameArg)
	}

	smcSecrets, ok := c.secrets[smcPubKey(smcKey)]
	if !ok {
		return xerrors.Errorf(errorKeyNotFoundInSmcs, smcKey)
	}

	_, found := smcSecrets[string(name)]

	if !found {
		return xerrors.Errorf(
			"'%s' was not found among the secrets of the smc (%v)",
			name, smcKey)
	}

	logs, err := getAuditLogs(snap, name)
	if err != nil {
		return xerrors.Errorf("failed to get audit logs for '%s': %v", name, err)
	}

	fmt.Fprintf(c.printer, "Audit logs for secret '%s':\n", name)

	for _, log := range logs {
		pubKey, err := getSecretAccess(snap, log)
		if err != nil {
			return xerrors.Errorf("failed to get public key for access token '%s': %v", log, err)
		}
		fmt.Fprintf(c.printer, "%x\n", pubKey)
	}

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

//
// Utility functions
//

func getSmcRoster(snap store.Snapshot, key []byte) ([]byte, error) {
	k := prefixed.NewPrefixedKey([]byte(PrefixSmcRosterKeys), key)
	roster, err := snap.Get(k)
	if err != nil {
		return nil, err
	}

	return roster, nil
}

func setSmcRoster(snap store.Snapshot, key []byte, roster []byte) error {
	k := prefixed.NewPrefixedKey([]byte(PrefixSmcRosterKeys), key)
	err := snap.Set(k, roster)
	if err != nil {
		return err
	}

	dela.Logger.Info().Str("contract", ContractName).
		Msgf("setting %x=%s", key, roster)

	return nil
}

func getSecret(snap store.Snapshot, key []byte) ([]byte, error) {
	k := prefixed.NewPrefixedKey([]byte(PrefixSecretKeys), key)
	secret, err := snap.Get(k)
	if secret == nil {
		if err != nil {
			// InMemorySnapshot version
			return nil, err
		}

		// BBolt store version
		return nil, xerrors.Errorf("couldn't find key")
	}

	return secret, nil
}

func setSecret(snap store.Snapshot, key []byte, secret []byte) error {
	k := prefixed.NewPrefixedKey([]byte(PrefixSecretKeys), key)
	err := snap.Set(k, secret)
	if err != nil {
		return err
	}

	dela.Logger.Info().
		Str("contract", ContractName).
		Msgf("setting secret %x=%s", key, secret)

	return nil
}

func deleteSecret(snap store.Snapshot, key []byte) error {
	k := prefixed.NewPrefixedKey([]byte(PrefixSecretKeys), key)
	err := snap.Delete(k)
	if err != nil {
		return xerrors.Errorf("failed to delete from snapshot: %v", err)
	}

	dela.Logger.Info().
		Str("contract", ContractName).
		Msgf("deleting secret %x", key)

	return nil
}

func computeAccessToken(smcKey []byte, secret []byte, clientPubKey []byte) []byte {
	h := crypto.NewHashFactory(crypto.Sha256).New()
	h.Write(smcKey)
	h.Write(secret)
	h.Write(clientPubKey)
	return h.Sum(nil)
}

func hasSecretAccess(snap store.Snapshot, accessToken []byte) bool {
	_, err := getSecretAccess(snap, accessToken)
	return err != nil
}

func getSecretAccess(snap store.Snapshot, accessToken []byte) ([]byte, error) {
	k := prefixed.NewPrefixedKey([]byte(PrefixAccessKeys), accessToken)
	pubKey, err := snap.Get(k)
	if err != nil {
		return nil, xerrors.Errorf(
			"failed to get access token '%v': %v", accessToken, err)
	}

	return pubKey, nil
}

func setSecretAccess(snap store.Snapshot, accessToken []byte, pubKey []byte) error {
	k := prefixed.NewPrefixedKey([]byte(PrefixAccessKeys), accessToken)
	err := snap.Set(k, pubKey)
	if err != nil {
		return xerrors.Errorf(
			"failed to give secret '%v' access to '%v': %v", accessToken, pubKey, err)
	}

	dela.Logger.Info().
		Str("contract", ContractName).
		Msgf("setting secret access %x=%s", accessToken, pubKey)

	return nil
}

func insertAuditLog(snap store.Snapshot, name []byte, accessToken []byte) error {
	k := prefixed.NewPrefixedKey([]byte(PrefixListKeys), name)
	log, err := snap.Get(k)
	if err != nil {
		log = make([]byte, 0, accessTokenNbBytes)
	}

	// log contains: [<accessToken1>, <accessToken2>, ...] , [][]byte
	log = append(log, accessToken...)

	err = snap.Set(k, log)
	if err != nil {
		return xerrors.Errorf(
			"failed to insert audit log for secret '%v': %v", name, err)
	}

	dela.Logger.Info().
		Str("contract", ContractName).
		Msgf("appending audit log %x=[%s]", name, accessToken)

	return nil
}

func getAuditLogs(snap store.Snapshot, name []byte) ([][]byte, error) {
	k := prefixed.NewPrefixedKey([]byte(PrefixListKeys), name)
	log, err := snap.Get(k)
	if err != nil {
		return [][]byte{}, err
	}

	return decodeAuditLog(log)
}

func decodeAuditLog(log []byte) ([][]byte, error) {
	// log contains: [<accessToken1>, <accessToken2>, ...] , [][]byte
	// we need to split it in chunks of accessTokenNbBytes bytes (sha256)
	// see computeAccessToken() above
	if len(log)%accessTokenNbBytes != 0 {
		return nil, xerrors.Errorf("invalid audit log length: %v", len(log))
	}

	res := make([][]byte, 0, len(log)/accessTokenNbBytes)
	for i := 0; i < len(log); i += accessTokenNbBytes {
		res = append(res, log[i:i+accessTokenNbBytes])
	}

	return res, nil
}
