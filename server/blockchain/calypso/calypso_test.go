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
	"go.dedis.ch/dela/core/store/prefixed"
	"go.dedis.ch/dela/core/txn"
	"go.dedis.ch/dela/core/txn/signed"
	"go.dedis.ch/dela/testing/fake"
)

func Test_secretSet_addSecret(t *testing.T) {
	s := secretSet{}
	secret := "some secret"

	s.addSecret(secret)

	require.Equal(t, 1, len(s))
	require.Equal(t, struct{}{}, s[secret])
}

func Test_NewCreds(t *testing.T) {
	creds := NewCreds()

	require.NotNilf(t, creds, "creds should not be nil")
}

func Test_RegisterContract(_ *testing.T) {
	RegisterContract(native.NewExecution(), Contract{})
}

func Test_NewContract(t *testing.T) {
	contract := NewContract(fakeAccess{})

	require.NotNilf(t, contract, "contract should not be nil")
	require.NotNilf(t, contract.secrets, "secrets should not be nil")
	require.NotNilf(t, contract.printer, "printer should not be nil")
	require.NotNilf(t, contract.cmd, "cmd should not be nil")
	require.NotNilf(t, contract.access, "access should not be nil")
}

func TestContract_UID(t *testing.T) {
	contract := NewContract(fakeAccess{})

	require.NotNilf(t, contract, "contract should not be nil")
	require.Equal(t, ContractUID, contract.UID())
}

func TestExecuteFailing(t *testing.T) {
	contract := NewContract(fakeAccess{err: fake.GetError()})

	err := contract.Execute(fakeStore{}, makeStep(t))
	require.EqualError(t, err,
		"identity not authorized: fake.PublicKey ("+fake.GetError().Error()+")")

	contract = NewContract(fakeAccess{})
	err = contract.Execute(fakeStore{}, makeStep(t))
	require.EqualError(t, err, "'calypso:command' not found in tx arg")

	contract.cmd = fakeCmd{err: fake.GetError()}

	err = contract.Execute(fakeStore{}, makeStep(t, CmdArg, "ADVERTISE_SMC"))
	require.EqualError(t, err, fake.Err("failed to ADVERTISE_SMC"))

	err = contract.Execute(fakeStore{}, makeStep(t, CmdArg, "DELETE_SMC"))
	require.EqualError(t, err, fake.Err("failed to DELETE_SMC"))

	err = contract.Execute(fakeStore{}, makeStep(t, CmdArg, "LIST_SMC"))
	require.EqualError(t, err, fake.Err("failed to LIST_SMC"))

	err = contract.Execute(fakeStore{}, makeStep(t, CmdArg, "CREATE_SECRET"))
	require.EqualError(t, err, fake.Err("failed to CREATE_SECRET"))

	err = contract.Execute(fakeStore{}, makeStep(t, CmdArg, "LIST_SECRETS"))
	require.EqualError(t, err, fake.Err("failed to LIST_SECRETS"))

	err = contract.Execute(fakeStore{}, makeStep(t, CmdArg, "REVEAL_SECRET"))
	require.EqualError(t, err, fake.Err("failed to REVEAL_SECRET"))

	err = contract.Execute(fakeStore{}, makeStep(t, CmdArg, "LIST_AUDIT_LOG"))
	require.EqualError(t, err, fake.Err("failed to LIST_AUDIT_LOG"))

	err = contract.Execute(fakeStore{}, makeStep(t, CmdArg, "fake"))
	require.EqualError(t, err, "unknown command: fake")
}

func TestExecutePassing(t *testing.T) {
	contract := NewContract(fakeAccess{})
	contract.cmd = fakeCmd{}

	err := contract.Execute(fakeStore{}, makeStep(t, CmdArg, "ADVERTISE_SMC"))
	require.NoError(t, err)

	err = contract.Execute(fakeStore{}, makeStep(t, CmdArg, "DELETE_SMC"))
	require.NoError(t, err)

	err = contract.Execute(fakeStore{}, makeStep(t, CmdArg, "LIST_SMC"))
	require.NoError(t, err)

	err = contract.Execute(fakeStore{}, makeStep(t, CmdArg, "CREATE_SECRET"))
	require.NoError(t, err)

	err = contract.Execute(fakeStore{}, makeStep(t, CmdArg, "LIST_SECRETS"))
	require.NoError(t, err)

	err = contract.Execute(fakeStore{}, makeStep(t, CmdArg, "REVEAL_SECRET"))
	require.NoError(t, err)

	err = contract.Execute(fakeStore{}, makeStep(t, CmdArg, "LIST_AUDIT_LOG"))
	require.NoError(t, err)

	err = contract.Execute(fakeStore{}, makeStep(t, CmdArg, "fake"))
	require.EqualError(t, err, "unknown command: fake")
}

func TestCommand_AdvertiseSmc(t *testing.T) {
	contract := NewContract(fakeAccess{})

	cmd := calypsoCommand{
		Contract: &contract,
	}

	keyString := "dummy"
	keySmc := smcPubKey(keyString)
	keyBytes := []byte(keyString)

	err := cmd.advertiseSmc(fake.NewSnapshot(), makeStep(t))
	require.EqualError(t, err, "'calypso:smc_key' not found in tx arg")

	err = cmd.advertiseSmc(fake.NewSnapshot(), makeStep(t, SmcPublicKeyArg, keyString))
	require.EqualError(t, err, "'calypso:smc_roster' not found in tx arg")

	err = cmd.advertiseSmc(fake.NewBadSnapshot(),
		makeStep(t, SmcPublicKeyArg, keyString, RosterArg, "node:12345"))
	require.EqualError(t, err, fake.Err("failed to set roster"))

	err = cmd.advertiseSmc(fake.NewBadSnapshot(),
		makeStep(t, SmcPublicKeyArg, keyString, RosterArg, ","))
	require.ErrorContains(t, err, "invalid node '' in roster")

	err = cmd.advertiseSmc(fake.NewBadSnapshot(),
		makeStep(t, SmcPublicKeyArg, keyString, RosterArg, "abcd"))
	require.ErrorContains(t, err, "invalid node 'abcd' in roster")

	snapshot := fake.NewSnapshot()

	_, found := contract.secrets[keySmc]
	require.False(t, found)

	err = cmd.advertiseSmc(snapshot,
		makeStep(t, SmcPublicKeyArg, keyString, RosterArg, "node:12345"))
	require.NoError(t, err)

	_, found = contract.secrets[keySmc]
	require.True(t, found)

	k := prefixed.NewPrefixedKey([]byte(PrefixSmcRosterKeys), keyBytes)
	res, err := snapshot.Get(k)
	require.NoError(t, err)
	require.Equal(t, "node:12345", string(res))
}

func TestCommand_DeleteSmc(t *testing.T) {
	contract := NewContract(fakeAccess{})

	cmd := calypsoCommand{
		Contract: &contract,
	}

	keyString := "dummy"
	keySmc := smcPubKey("dummy")
	keyBytes := []byte(keyString)
	keyHex := hex.EncodeToString(keyBytes)

	err := cmd.deleteSmc(fake.NewSnapshot(), makeStep(t))
	require.EqualError(t, err, "'calypso:smc_key' not found in tx arg")

	err = cmd.deleteSmc(fake.NewBadSnapshot(), makeStep(t, SmcPublicKeyArg, keyString))
	require.EqualError(t, err, fake.Err("failed to delete SMC with public key '"+keyHex+"'"))

	snap := fake.NewSnapshot()
	err = snap.Set(keyBytes, []byte("localhost:12345"))
	require.NoError(t, err)
	contract.secrets[keySmc] = secretSet{}

	err = cmd.deleteSmc(snap, makeStep(t, SmcPublicKeyArg, keyString))
	require.NoError(t, err)

	k := prefixed.NewPrefixedKey([]byte(PrefixSmcRosterKeys), keyBytes)
	res, err := snap.Get(k)
	require.Nil(t, err) // == key not found
	require.Nil(t, res)

	_, found := contract.secrets[keySmc]
	require.False(t, found)
}

func TestCommand_ListSmc(t *testing.T) {
	contract := NewContract(fakeAccess{})

	key1String := "key1"
	key1Smc := smcPubKey(key1String)
	key1Bytes := []byte(key1String)
	roster1 := "localhost:12345"

	key2String := "key2"
	key2Smc := smcPubKey(key2String)
	key2Bytes := []byte(key2String)
	roster2 := "localhost:12345,remote:54321"

	contract.secrets[key1Smc] = secretSet{}
	contract.secrets[key2Smc] = secretSet{}

	buf := &bytes.Buffer{}
	contract.printer = buf

	cmd := calypsoCommand{
		Contract: &contract,
	}

	snap := fake.NewSnapshot()

	k := prefixed.NewPrefixedKey([]byte(PrefixSmcRosterKeys), key1Bytes)
	err := snap.Set(k, []byte(roster1))
	require.NoError(t, err)

	k = prefixed.NewPrefixedKey([]byte(PrefixSmcRosterKeys), key2Bytes)
	err = snap.Set(k, []byte(roster2))
	require.NoError(t, err)

	err = cmd.listSmc(snap)
	require.NoError(t, err)

	require.Equal(t, fmt.Sprintf("%x=%v,%x=%v", key1String, roster1, key2String, roster2),
		buf.String())

	err = cmd.listSmc(fake.NewBadSnapshot())
	// we can't assume an order from the map
	require.Regexp(t, "^failed to get key", err.Error())
}

func TestCommand_CreateSecret_BadSnapshot(t *testing.T) {
	// Arrange
	contract := NewContract(fakeAccess{})

	cmd := calypsoCommand{
		Contract: &contract,
	}

	badSnap := fake.NewBadSnapshot()
	badSnap.ErrRead = nil  // temporarily disable errors
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

func TestCommand_CreateSecret_AlreadyExists(t *testing.T) {
	// Arrange
	contract := NewContract(fakeAccess{})

	cmd := calypsoCommand{
		Contract: &contract,
	}

	badSnap := fake.NewSnapshot()

	err := cmd.advertiseSmc(badSnap, makeStep(t, SmcPublicKeyArg, "dummy", RosterArg, "node:12345"))
	require.NoError(t, err)

	err = cmd.createSecret(badSnap,
		makeStep(t, SmcPublicKeyArg, "dummy", SecretNameArg, "name", SecretArg, "value"))
	require.NoError(t, err)

	// Act
	err = cmd.createSecret(badSnap,
		makeStep(t, SmcPublicKeyArg, "dummy", SecretNameArg, "name", SecretArg, "other_value"))

	// Assert
	require.EqualError(t, err, "a secret named 'name' already exists")
}

func TestCommand_CreateSecret_Succeeds(t *testing.T) {
	// Arrange
	contract := NewContract(fakeAccess{})

	cmd := calypsoCommand{
		Contract: &contract,
	}

	snap := fake.NewSnapshot()
	err := cmd.advertiseSmc(snap, makeStep(t, SmcPublicKeyArg, "dummy", RosterArg, "node:12345"))
	require.NoError(t, err)

	// Verify pre-conditions
	dummy, found := contract.secrets["dummy"]
	require.True(t, found)
	require.Equal(t, 0, len(dummy))

	// Act
	err = cmd.createSecret(snap,
		makeStep(t, SmcPublicKeyArg, "dummy", SecretNameArg, "my_secret", SecretArg, "my_value"))

	// Assert
	require.NoError(t, err)

	dummy = contract.secrets[smcPubKey("dummy")]
	require.Equal(t, 1, len(dummy))
	require.Equal(t, dummy["my_value"], struct{}{})

	k := prefixed.NewPrefixedKey([]byte(PrefixSecretKeys), []byte("my_secret"))
	res, err := snap.Get(k)
	require.NoError(t, err)
	require.Equal(t, "my_value", string(res))
}

func TestCommand_CreateSecret_InvalidInputs(t *testing.T) {
	contract := NewContract(fakeAccess{})

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

	contract := NewContract(fakeAccess{})

	buf := &bytes.Buffer{}
	contract.printer = buf

	cmd := calypsoCommand{
		Contract: &contract,
	}

	snap := fake.NewSnapshot()

	err := cmd.advertiseSmc(snap,
		makeStep(t, SmcPublicKeyArg, "dummy", RosterArg, "node:12345"))
	require.NoError(t, err)

	err = cmd.advertiseSmc(snap,
		makeStep(t, SmcPublicKeyArg, "other", RosterArg, "node:32145"))
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
	_, found := contract.secrets["dummy"]
	require.True(t, found)

	require.Equal(t, 2, len(contract.secrets["dummy"]))

	_, found = contract.secrets["other"]
	require.True(t, found)

	require.Equal(t, 1, len(contract.secrets["other"]))

	// Act
	err = cmd.listSecrets(snap,
		makeStep(t, SmcPublicKeyArg, "dummy"))

	// Assert
	require.NoError(t, err)
	require.Equal(t, "name1=secret1,name2=secret2", buf.String())
}

func TestCommand_ListSecrets_EmptySmcKey(t *testing.T) {
	// Arrange
	contract := NewContract(fakeAccess{})

	buf := &bytes.Buffer{}
	contract.printer = buf

	cmd := calypsoCommand{
		Contract: &contract,
	}

	// Act
	err := cmd.listSecrets(fake.NewSnapshot(), makeStep(t, SmcPublicKeyArg, ""))

	// Assert
	require.Error(t, err, "'calypso:smc_key' not found in tx arg")
}

func TestCommand_ListSecrets_NonexistentSmc(t *testing.T) {
	// Arrange
	contract := NewContract(fakeAccess{})

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
	contract := NewContract(fakeAccess{})

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
	contract := NewContract(fakeAccess{})

	cmd := calypsoCommand{
		Contract: &contract,
	}

	snap := fake.NewSnapshot()
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
	smcSecrets, found := contract.secrets[smcKey]
	require.True(t, found)
	require.Equal(t, 1, len(smcSecrets))

	_, found = smcSecrets[secretName]
	require.True(t, found)

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

func TestCommand_RevealSecretFails(t *testing.T) {
	contract := NewContract(fakeAccess{})
	cmd := calypsoCommand{
		Contract: &contract,
	}

	snap := fake.NewSnapshot()

	const (
		smcKey      = "my_smc_key"
		secretName  = "my_secret"
		secretValue = "my_value"
		pubKey      = "my_pubkey"
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
	smcSecrets, found := contract.secrets[smcKey]
	require.True(t, found)
	require.Equal(t, 1, len(smcSecrets))

	_, found = smcSecrets[secretName]
	require.True(t, found)

	err = cmd.revealSecret(snap, makeStep(t,
		SmcPublicKeyArg, "",
		SecretNameArg, secretName,
		PubKeyArg, pubKey))
	require.EqualError(t, err, "'calypso:smc_key' not found in tx arg")

	err = cmd.revealSecret(snap, makeStep(t,
		SmcPublicKeyArg, smcKey,
		SecretNameArg, "",
		PubKeyArg, pubKey))
	require.EqualError(t, err, "'calypso:secret_name' not found in tx arg")

	err = cmd.revealSecret(snap, makeStep(t,
		SmcPublicKeyArg, smcKey,
		SecretNameArg, secretName,
		PubKeyArg, ""))
	require.EqualError(t, err, "'calypso:pub_key' not found in tx arg")

	err = cmd.revealSecret(snap, makeStep(t,
		SmcPublicKeyArg, "another_smc_key",
		SecretNameArg, secretName,
		PubKeyArg, pubKey))
	require.EqualError(t, err, "'another_smc_key' was not found among the SMCs")

	err = cmd.revealSecret(snap, makeStep(t,
		SmcPublicKeyArg, smcKey,
		SecretNameArg, "another_name",
		PubKeyArg, pubKey))
	require.EqualError(t, err,
		"'another_name' was not found among the secrets of the smc (my_smc_key)")
}

func TestCommand_ListAuditLogs_Succeeds(t *testing.T) {

	// Arrange
	contract := NewContract(fakeAccess{})

	buf := &bytes.Buffer{}
	contract.printer = buf

	cmd := calypsoCommand{
		Contract: &contract,
	}

	snap := fake.NewSnapshot()

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
	_, found := contract.secrets[smcKey]
	require.True(t, found)

	smcSecrets := contract.secrets[smcKey]
	require.Equal(t, 1, len(smcSecrets))
	require.Equal(t, smcSecrets[secretName], struct{}{})

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

func TestRegisterContract(_ *testing.T) {
	RegisterContract(native.NewExecution(), Contract{})
}

func Test_validateRosterUpdate(t *testing.T) {
	type args struct {
		oldRoster []byte
		newRoster []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "validateNew",
			args: args{
				oldRoster: []byte(""),
				newRoster: []byte("node:12345"),
			},
			wantErr: true,
		},
		{
			name: "validateAdd",
			args: args{
				oldRoster: []byte("node:12345"),
				newRoster: []byte("node:12345,other:54321"),
			},
			wantErr: false,
		},
		{
			name: "validateNoChange",
			args: args{
				oldRoster: []byte("node:12345"),
				newRoster: []byte("node:12345"),
			},
			wantErr: false,
		},
		{
			name: "validateRemove",
			args: args{
				oldRoster: []byte("node:12345,other:54321"),
				newRoster: []byte("node:12345"),
			},
			wantErr: true,
		},
		{
			name: "validateNoOverlap",
			args: args{
				oldRoster: []byte("node:12345"),
				newRoster: []byte("other:54321"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateRosterUpdate(tt.args.oldRoster,
				tt.args.newRoster); (err != nil) != tt.wantErr {
				t.Errorf("validateRosterUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
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

func (s fakeStore) Get(_ []byte) ([]byte, error) {
	return nil, nil
}

func (s fakeStore) Set(_, _ []byte) error {
	return nil
}

type fakeCmd struct {
	err error
}

func (c fakeCmd) advertiseSmc(_ store.Snapshot, _ execution.Step) error {
	return c.err
}

func (c fakeCmd) deleteSmc(_ store.Snapshot, _ execution.Step) error {
	return c.err
}

func (c fakeCmd) listSmc(_ store.Snapshot) error {
	return c.err
}

func (c fakeCmd) createSecret(_ store.Snapshot, _ execution.Step) error {
	return c.err
}

func (c fakeCmd) listSecrets(_ store.Snapshot, _ execution.Step) error {
	return c.err
}

func (c fakeCmd) revealSecret(_ store.Snapshot, _ execution.Step) error {
	return c.err
}

func (c fakeCmd) listAuditLogs(_ store.Snapshot, _ execution.Step) error {
	return c.err
}
