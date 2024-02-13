package controller

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"go.dedis.ch/dela"
	"go.dedis.ch/dela/cli/node"
	"go.dedis.ch/dela/dkg"
	"go.dedis.ch/kyber/v3/util/key"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/suites"
	"golang.org/x/xerrors"
)

// suite is the Kyber suite for Pedersen.
var suite = suites.MustFind("Ed25519")

const separator = ":"
const malformedEncoded = "malformed encoded: %s"
const keyFileName = "key.pair"

// createKeyPairAction is an action to create a key pair
//
// - implements node.ActionTemplate
type createKeyPairAction struct{}

func (c createKeyPairAction) Execute(_ node.Context) error {
	kp := key.NewKeyPair(suites.MustFind("Ed25519"))

	privk, err := kp.Private.MarshalBinary()
	if err != nil {
		return xerrors.Errorf("failed to marshal private key: %v", err)
	}

	pubk, err := kp.Public.MarshalBinary()
	if err != nil {
		return xerrors.Errorf("failed to marshal public key: %v", err)
	}

	keyFile, err := os.Create(keyFileName)
	if err != nil {
		return xerrors.Errorf("failed to create key file: %v", err)
	}
	defer func() {
		if err := keyFile.Close(); err != nil {
			panic(err)
		}
	}()

	fmt.Fprintf(keyFile, "%v%v%v\n",
		hex.EncodeToString(privk),
		separator,
		hex.EncodeToString(pubk))

	return nil
}

// revealAction is an action to reveal a message
//
// - implements node.ActionTemplate
type revealAction struct{}

func (r revealAction) Execute(ctx node.Context) error {
	xhatString := ctx.Flags.String("xhatenc")
	xhatenc, err := decodePublicKey(xhatString)
	if err != nil {
		return xerrors.Errorf("failed to reveal: %v", err)
	}

	dkgpubString := ctx.Flags.String("dkgpub")
	var dkgpubk kyber.Point
	if dkgpubString != "" {
		dkgpubk, err = decodePublicKey(dkgpubString)
		if err != nil {
			return xerrors.Errorf("failed to decode public key str: %v", err)
		}
	} else {
		var actor dkg.Actor
		err := ctx.Injector.Resolve(&actor)
		if err != nil {
			return xerrors.Errorf("failed to resolve DKG actor: %v", err)
		}

		dkgpubk, err = actor.GetPublicKey()
		if err != nil {
			return xerrors.Errorf("failed retrieving DKG public key: %v", err)
		}
	}

	privkString := ctx.Flags.String("privk")
	privateKey, err := decodePrivateKey(privkString)
	if err != nil {
		return xerrors.Errorf("failed to decode private key str: %v", err)
	}

	encrypted := ctx.Flags.String("encrypted")
	_, cs, err := decodeEncrypted(encrypted)
	if err != nil {
		return xerrors.Errorf("failed to decode encrypted str: %v", err)
	}

	msg, err := reveal(xhatenc, dkgpubk, privateKey, cs)
	if err != nil {
		fmt.Printf("couldn't reveal message. %v", err)
		return err
	}
	fmt.Print(hex.EncodeToString(msg))

	return nil
}

// -----------------------------------------------------------------------------
// helper functions
func decodePrivateKey(str string) (kyber.Scalar, error) {
	pkbuff, err := hex.DecodeString(str)
	if err != nil {
		return nil, xerrors.Errorf(malformedEncoded, str)
	}

	pk := suite.Scalar()
	err = pk.UnmarshalBinary(pkbuff)
	if err != nil {
		return nil, xerrors.Errorf("failed to unmarshal pk: %v", err)
	}

	return pk, nil
}

func decodePublicKey(str string) (kyber.Point, error) {
	pkbuff, err := hex.DecodeString(str)
	if err != nil {
		return nil, xerrors.Errorf(malformedEncoded, str)
	}

	pk := suite.Point()
	err = pk.UnmarshalBinary(pkbuff)
	if err != nil {
		return nil, xerrors.Errorf("failed to unmarshal pk: %v", err)
	}

	return pk, nil
}

func decodeEncrypted(str string) (kyber.Point, []kyber.Point, error) {
	parts := strings.Split(str, separator)
	if len(parts) < 2 {
		return nil, nil, xerrors.Errorf(malformedEncoded, str)
	}

	// Decode K
	kbuff, err := hex.DecodeString(parts[0])
	if err != nil {
		return nil, nil, xerrors.Errorf("failed to decode k point: %v", err)
	}

	k := suite.Point()
	err = k.UnmarshalBinary(kbuff)
	if err != nil {
		return nil, nil, xerrors.Errorf("failed to unmarshal k point: %v", err)
	}

	// Decode Cs
	cs := make([]kyber.Point, 0, len(parts)-1)
	for _, p := range parts[1:] {
		cbuff, err := hex.DecodeString(p)
		if err != nil {
			return nil, nil, xerrors.Errorf("failed to decode c point: %v", err)
		}

		c := suite.Point()
		err = c.UnmarshalBinary(cbuff)
		if err != nil {
			return nil, nil, xerrors.Errorf("failed to unmarshal c point: %v", err)
		}

		cs = append(cs, c)
	}

	dela.Logger.Debug().Msgf("Decoded K: %v and Cs: %v", k, cs)

	return k, cs, nil
}
