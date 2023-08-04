package calypso

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.dedis.ch/dela/core/access"
	"go.dedis.ch/dela/core/execution"
	"go.dedis.ch/dela/core/execution/native"
	"go.dedis.ch/dela/core/store"
	"go.dedis.ch/dela/core/txn"
	"go.dedis.ch/dela/core/txn/signed"
	"go.dedis.ch/dela/testing/fake"
	"golang.org/x/xerrors"
)

func TestExecute(t *testing.T) {
	contract := NewContract([]byte{}, fakeAccess{err: fake.GetError()})

	err := contract.Execute(fakeStore{}, makeStep(t))
	require.EqualError(t, err,
		"identity not authorized: fake.PublicKey ("+fake.GetError().Error()+")")

	contract = NewContract([]byte{}, fakeAccess{})
	err = contract.Execute(fakeStore{}, makeStep(t))
	require.EqualError(t, err, "'calypso:command' not found in tx arg")

	contract.cmd = fakeCmd{err: fake.GetError()}

	err = contract.Execute(fakeStore{}, makeStep(t, CmdArg, "ADVERTISE_SMC"))
	require.EqualError(t, err, fake.Err("failed to ADVERTISE_SMC"))

	err = contract.Execute(fakeStore{}, makeStep(t, CmdArg, "DELETE_SMC"))
	require.EqualError(t, err, fake.Err("failed to DELETE_SMC"))

	err = contract.Execute(fakeStore{}, makeStep(t, CmdArg, "LIST_SMC"))
	require.EqualError(t, err, fake.Err("failed to LIST_SMC"))

	err = contract.Execute(fakeStore{}, makeStep(t, CmdArg, "fake"))
	require.EqualError(t, err, "unknown command: fake")

	contract.cmd = fakeCmd{}
	err = contract.Execute(fakeStore{}, makeStep(t, CmdArg, "ADVERTISE_SMC"))
	require.NoError(t, err)
}

func TestCommand_AdvertiseSmc(t *testing.T) {
	contract := NewContract([]byte{}, fakeAccess{})

	cmd := calypsoCommand{
		Contract: &contract,
	}

	err := cmd.advertiseSmc(fake.NewSnapshot(), makeStep(t))
	require.EqualError(t, err, "'calypso:smc_key' not found in tx arg")

	err = cmd.advertiseSmc(fake.NewSnapshot(), makeStep(t, SmcPublicKeyArg, "dummy"))
	require.EqualError(t, err, "'calypso:smc_roster' not found in tx arg")

	err = cmd.advertiseSmc(fake.NewBadSnapshot(),
		makeStep(t, SmcPublicKeyArg, "dummy", RosterArg, "node:12345"))
	require.EqualError(t, err, fake.Err("failed to set roster"))

	err = cmd.advertiseSmc(fake.NewBadSnapshot(),
		makeStep(t, SmcPublicKeyArg, "dummy", RosterArg, ","))
	require.ErrorContains(t, err, "invalid node '' in roster")

	err = cmd.advertiseSmc(fake.NewBadSnapshot(),
		makeStep(t, SmcPublicKeyArg, "dummy", RosterArg, "abcd"))
	require.ErrorContains(t, err, "invalid node 'abcd' in roster")

	snap := fake.NewSnapshot()

	_, found := contract.index["dummy"]
	require.False(t, found)

	_, found = contract.secrets["dummy"]
	require.False(t, found)

	err = cmd.advertiseSmc(snap, makeStep(t, SmcPublicKeyArg, "dummy", RosterArg, "node:12345"))
	require.NoError(t, err)

	_, found = contract.index["dummy"]
	require.True(t, found)

	_, found = contract.secrets["dummy"]
	require.True(t, found)

	res, err := snap.Get([]byte("dummy"))
	require.NoError(t, err)
	require.Equal(t, "node:12345", string(res))
}

func TestCommand_DeleteSmc(t *testing.T) {
	contract := NewContract([]byte{}, fakeAccess{})

	cmd := calypsoCommand{
		Contract: &contract,
	}

	key := []byte("dummy")
	keyHex := hex.EncodeToString(key)
	keyStr := string(key)

	err := cmd.deleteSmc(fake.NewSnapshot(), makeStep(t))
	require.EqualError(t, err, "'calypso:smc_key' not found in tx arg")

	err = cmd.deleteSmc(fake.NewBadSnapshot(), makeStep(t, SmcPublicKeyArg, keyStr))
	require.EqualError(t, err, fake.Err("failed to deleteSmc key '"+keyHex+"'"))

	snap := fake.NewSnapshot()
	snap.Set(key, []byte("localhost:12345"))
	contract.index[keyStr] = struct{}{}

	err = cmd.deleteSmc(snap, makeStep(t, SmcPublicKeyArg, keyStr))
	require.NoError(t, err)

	res, err := snap.Get(key)
	require.Nil(t, err)
	require.Nil(t, res)

	_, found := contract.index[keyStr]
	require.False(t, found)
}

func TestCommand_ListSmc(t *testing.T) {
	contract := NewContract([]byte{}, fakeAccess{})

	key1 := "key1"
	roster1 := "localhost:12345"

	key2 := "key2"
	roster2 := "localhost:12345,remote:54321"

	contract.index[key1] = struct{}{}
	contract.index[key2] = struct{}{}

	buf := &bytes.Buffer{}
	contract.printer = buf

	cmd := calypsoCommand{
		Contract: &contract,
	}

	snap := fake.NewSnapshot()
	snap.Set([]byte(key1), []byte(roster1))
	snap.Set([]byte(key2), []byte(roster2))

	err := cmd.listSmc(snap)
	require.NoError(t, err)

	require.Equal(t, fmt.Sprintf("%x=%v,%x=%v", key1, roster1, key2, roster2), buf.String())

	err = cmd.listSmc(fake.NewBadSnapshot())
	// we can't assume an order from the map
	require.Regexp(t, "^failed to get key", err.Error())
}

func TestCommand_CreateSecret_BadSnapshot(t *testing.T) {
	// Arrange
	contract := NewContract([]byte{}, fakeAccess{})

	cmd := calypsoCommand{
		Contract: &contract,
	}

	badSnap := fake.NewBadSnapshot()
	badSnap.ErrWrite = nil // temporarily disable errors

	err := cmd.advertiseSmc(badSnap, makeStep(t, SmcPublicKeyArg, "dummy", RosterArg, "node:12345"))
	require.NoError(t, err)

	badSnap.ErrWrite = fake.GetError() // re-enable errors

	// Act
	err = cmd.createSecret(badSnap,
		makeStep(t, SmcPublicKeyArg, "dummy", SecretNameArg, "name", SecretArg, "value"))

	// Assert
	require.EqualError(t, err, fake.Err("failed to set secret"))
}

func TestCommand_CreateSecret_Succeeds(t *testing.T) {
	// Arrange
	contract := NewContract([]byte{}, fakeAccess{})

	cmd := calypsoCommand{
		Contract: &contract,
	}

	snap := fake.NewSnapshot()
	err := cmd.advertiseSmc(snap, makeStep(t, SmcPublicKeyArg, "dummy", RosterArg, "node:12345"))
	require.NoError(t, err)

	// Verify pre-conditions
	_, found := contract.index["dummy"]
	require.True(t, found)

	_, found = contract.secrets["dummy"]
	require.True(t, found)

	dummy := contract.secrets["dummy"]
	require.Equal(t, 0, len(dummy))

	// Act
	err = cmd.createSecret(snap,
		makeStep(t, SmcPublicKeyArg, "dummy", SecretNameArg, "my_secret", SecretArg, "my_value"))

	// Assert
	require.NoError(t, err)

	dummy = contract.secrets["dummy"]
	require.Equal(t, 1, len(dummy))
	require.Equal(t, "my_secret", string(dummy[0]))

	res, err := snap.Get([]byte("my_secret"))
	require.NoError(t, err)
	require.Equal(t, "my_value", string(res))
}

func TestCommand_CreateSecret_InvalidInputs(t *testing.T) {
	contract := NewContract([]byte{}, fakeAccess{})

	cmd := calypsoCommand{
		Contract: &contract,
	}

	err := cmd.createSecret(fake.NewSnapshot(), makeStep(t))
	require.ErrorContains(t, err, "not found in tx arg")

	err = cmd.createSecret(fake.NewSnapshot(),
		makeStep(t, SecretNameArg, "name", SecretArg, "value"))
	require.EqualError(t, err, "'calypso:smc_key' not found in tx arg")

	err = cmd.createSecret(fake.NewSnapshot(),
		makeStep(t, SmcPublicKeyArg, "dummy", SecretNameArg, "name"))
	require.EqualError(t, err, "'calypso:secret_value' not found in tx arg")

	err = cmd.createSecret(fake.NewSnapshot(),
		makeStep(t, SmcPublicKeyArg, "dummy", SecretArg, "value"))
	require.EqualError(t, err, "'calypso:secret_name' not found in tx arg")

	err = cmd.createSecret(fake.NewSnapshot(),
		makeStep(t, SmcPublicKeyArg, "dummy", SecretNameArg, "name", SecretArg, ""))
	require.ErrorContains(t, err, "'calypso:secret_value' not found in tx arg")

	err = cmd.createSecret(fake.NewSnapshot(),
		makeStep(t, SmcPublicKeyArg, "dummy", SecretNameArg, "", SecretArg, "value"))
	require.ErrorContains(t, err, "'calypso:secret_name' not found in tx arg")

	err = cmd.createSecret(fake.NewSnapshot(),
		makeStep(t, SmcPublicKeyArg, "invalid", SecretNameArg, "n", SecretArg, "v"))
	require.ErrorContains(t, err, "'invalid' was not found among the SMCs")
}

func TestCommand_ListSecrets(t *testing.T) {

	// Arrange (2 SMCs, "dummy" has 2 secrets, "other" has 1)

	contract := NewContract([]byte{}, fakeAccess{})

	buf := &bytes.Buffer{}
	contract.printer = buf

	cmd := calypsoCommand{
		Contract: &contract,
	}

	snap := fake.NewSnapshot()

	err := cmd.advertiseSmc(snap, makeStep(t, SmcPublicKeyArg, "dummy", RosterArg, "node:12345"))
	require.NoError(t, err)

	err = cmd.advertiseSmc(snap, makeStep(t, SmcPublicKeyArg, "other", RosterArg, "node:32145"))
	require.NoError(t, err)

	err = cmd.createSecret(snap,
		makeStep(t, SmcPublicKeyArg, "dummy", SecretNameArg, "name1", SecretArg, "secret1"))
	require.NoError(t, err)

	err = cmd.createSecret(snap,
		makeStep(t, SmcPublicKeyArg, "dummy", SecretNameArg, "name2", SecretArg, "secret2"))
	require.NoError(t, err)

	err = cmd.createSecret(snap,
		makeStep(t, SmcPublicKeyArg, "other", SecretNameArg, "name3", SecretArg, "secret3"))
	require.NoError(t, err)

	// Verify pre-conditions
	_, found := contract.index["dummy"]
	require.True(t, found)

	_, found = contract.secrets["dummy"]
	require.True(t, found)

	require.Equal(t, 2, len(contract.secrets["dummy"]))

	_, found = contract.index["other"]
	require.True(t, found)

	_, found = contract.secrets["other"]
	require.True(t, found)

	require.Equal(t, 1, len(contract.secrets["other"]))

	// Act
	err = cmd.listSecrets(snap, makeStep(t, SmcPublicKeyArg, "dummy"))

	// Assert
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("name1=secret1,name2=secret2"), buf.String())
}

func TestCommand_ListSecrets_InexistentSmc(t *testing.T) {
	// Arrange
	contract := NewContract([]byte{}, fakeAccess{})

	buf := &bytes.Buffer{}
	contract.printer = buf

	cmd := calypsoCommand{
		Contract: &contract,
	}

	// Act
	err := cmd.listSecrets(fake.NewSnapshot(), makeStep(t, SmcPublicKeyArg, "noexist"))

	// Assert
	require.ErrorContains(t, err, "SMC not found: noexist")
}

func TestCommand_ListSecrets_InvalidSnapshot(t *testing.T) {
	// Arrange
	contract := NewContract([]byte{}, fakeAccess{})

	buf := &bytes.Buffer{}
	contract.printer = buf

	cmd := calypsoCommand{
		Contract: &contract,
	}

	snap := fake.NewBadSnapshot()
	snap.ErrWrite = nil
	snap.ErrRead = nil

	err := cmd.advertiseSmc(snap, makeStep(t, SmcPublicKeyArg, "dummy", RosterArg, "node:12345"))
	require.NoError(t, err)

	err = cmd.createSecret(snap,
		makeStep(t, SmcPublicKeyArg, "dummy", SecretNameArg, "name1", SecretArg, "secret1"))
	require.NoError(t, err)

	snap.ErrWrite = fake.GetError()
	snap.ErrRead = fake.GetError()

	// Act
	err = cmd.listSecrets(snap, makeStep(t, SmcPublicKeyArg, "dummy"))

	// Assert
	require.ErrorContains(t, err, "failed to get key")
}

func TestCommand_RevealSecret_Succeeds(t *testing.T) {

	// Arrange
	contract := NewContract([]byte{}, fakeAccess{})

	cmd := calypsoCommand{
		Contract: &contract,
	}

	snap := NewSnapshot()
	const (
		smcKey      = "my_smc_key"
		secretName  = "my_secret"
		secretValue = "my_value"
	)

	err := cmd.advertiseSmc(snap,
		makeStep(t,
			SmcPublicKeyArg, smcKey,
			RosterArg, "node:12345"))

	require.NoError(t, err)

	err = cmd.createSecret(snap,
		makeStep(t,
			SmcPublicKeyArg, smcKey,
			SecretNameArg, secretName,
			SecretArg, secretValue))
	require.NoError(t, err)

	// Verify pre-conditions
	_, found := contract.index[smcKey]
	require.True(t, found)

	_, found = contract.secrets[smcKey]
	require.True(t, found)

	smcSecrets := contract.secrets[smcKey]
	require.Equal(t, 1, len(smcSecrets))
	require.Equal(t, secretName, smcSecrets[0])

	// Act
	err = cmd.revealSecret(snap,
		makeStep(t,
			SmcPublicKeyArg, smcKey,
			SecretNameArg, secretName,
			PubKeyArg, "my_pubkey"))

	// Assert
	require.NoError(t, err)

	token := computeAccessToken([]byte(smcKey), []byte(secretValue), []byte("my_pubkey"))

	logs, err := getAuditLogs(snap, []byte(secretName))
	require.NoError(t, err)
	require.Len(t, logs, 1)
	require.Equal(t, logs[0], token)
}

func TestCommand_ListAuditLogs_Succeeds(t *testing.T) {

	// Arrange
	contract := NewContract([]byte{}, fakeAccess{})

	buf := &bytes.Buffer{}
	contract.printer = buf

	cmd := calypsoCommand{
		Contract: &contract,
	}

	snap := NewSnapshot()
	const (
		smcKey      = "my_smc_key"
		secretName  = "my_secret"
		secretValue = "my_value"
	)

	err := cmd.advertiseSmc(snap,
		makeStep(t,
			SmcPublicKeyArg, smcKey,
			RosterArg, "node:12345"))

	require.NoError(t, err)

	err = cmd.createSecret(snap,
		makeStep(t,
			SmcPublicKeyArg, smcKey,
			SecretNameArg, secretName,
			SecretArg, secretValue))
	require.NoError(t, err)

	// Verify pre-conditions
	_, found := contract.index[smcKey]
	require.True(t, found)

	_, found = contract.secrets[smcKey]
	require.True(t, found)

	smcSecrets := contract.secrets[smcKey]
	require.Equal(t, 1, len(smcSecrets))
	require.Equal(t, secretName, smcSecrets[0])

	err = cmd.revealSecret(snap,
		makeStep(t,
			SmcPublicKeyArg, smcKey,
			SecretNameArg, secretName,
			PubKeyArg, "my_pubkey"))
	require.NoError(t, err)

	// Act
	err = cmd.listAuditLogs(snap,
		makeStep(t,
			SmcPublicKeyArg, smcKey,
			SecretNameArg, secretName))

	// Assert
	require.NoError(t, err)

	require.Equal(t,
		fmt.Sprintf("Audit logs for secret '%v':\n", secretName)+
			fmt.Sprintf("%x\n", "my_pubkey"),
		buf.String())
}

func TestInfoLog(t *testing.T) {
	log := infoLog{}

	n, err := log.Write([]byte{0b0, 0b1})
	require.NoError(t, err)
	require.Equal(t, 2, n)
}

func TestRegisterContract(t *testing.T) {
	RegisterContract(native.NewExecution(), Contract{})
}

// -----------------------------------------------------------------------------
// Utility functions

func makeStep(t *testing.T, args ...string) execution.Step {
	return execution.Step{Current: makeTx(t, args...)}
}

func makeTx(t *testing.T, args ...string) txn.Transaction {
	options := []signed.TransactionOption{}
	for i := 0; i < len(args)-1; i += 2 {
		options = append(options, signed.WithArg(args[i], []byte(args[i+1])))
	}

	tx, err := signed.NewTransaction(0, fake.PublicKey{}, options...)
	require.NoError(t, err)

	return tx
}

type fakeAccess struct {
	access.Service

	err error
}

func (srvc fakeAccess) Match(store.Readable, access.Credential, ...access.Identity) error {
	return srvc.err
}

func (srvc fakeAccess) Grant(store.Snapshot, access.Credential, ...access.Identity) error {
	return srvc.err
}

type fakeStore struct {
	store.Snapshot
}

func (s fakeStore) Get(key []byte) ([]byte, error) {
	return nil, nil
}

func (s fakeStore) Set(key, value []byte) error {
	return nil
}

type fakeCmd struct {
	err error
}

func (c fakeCmd) advertiseSmc(snap store.Snapshot, step execution.Step) error {
	return c.err
}

func (c fakeCmd) deleteSmc(snap store.Snapshot, step execution.Step) error {
	return c.err
}

func (c fakeCmd) listSmc(snap store.Snapshot) error {
	return c.err
}

func (c fakeCmd) createSecret(snap store.Snapshot, step execution.Step) error {
	return c.err
}

func (c fakeCmd) listSecrets(snap store.Snapshot, step execution.Step) error {
	return c.err
}

func (c fakeCmd) revealSecret(snap store.Snapshot, step execution.Step) error {
	return c.err
}

func (c fakeCmd) listAuditLogs(snap store.Snapshot, step execution.Step) error {
	return c.err
}

// -----------------------------------------------------------------------------

// InMemorySnapshot is a fake, but realistic implementation of a store snapshot.
// TODO: should we update the fake.InMemorySnapshot ?
// - implements store.Snapshot
type InMemorySnapshot struct {
	store.Snapshot
	values map[string][]byte
}

// NewSnapshot creates a new empty snapshot.
func NewSnapshot() *InMemorySnapshot {
	return &InMemorySnapshot{
		values: make(map[string][]byte),
	}
}

// Get implements store.Snapshot.
func (snap *InMemorySnapshot) Get(key []byte) ([]byte, error) {
	value, found := snap.values[string(key)]
	if found {
		return value, nil
	}
	return nil, xerrors.Errorf("key not found: %s", key)
}

// Set implements store.Snapshot.
func (snap *InMemorySnapshot) Set(key, value []byte) error {
	snap.values[string(key)] = value
	return nil
}

// Delete implements store.Snapshot.
func (snap *InMemorySnapshot) Delete(key []byte) error {
	_, found := snap.values[string(key)]
	if !found {
		// is this behaviour correct or should it be ignored ?
		return xerrors.Errorf("key not found: %s", key)
	}

	delete(snap.values, string(key))
	return nil
}
