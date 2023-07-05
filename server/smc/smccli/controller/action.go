package controller

import (
	"encoding/hex"
	"fmt"
	"go.dedis.ch/kyber/v3/util/key"
	"strings"

	"go.dedis.ch/dela/cli/node"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/suites"
	"golang.org/x/xerrors"
)

// suite is the Kyber suite for Pedersen.
var suite = suites.MustFind("Ed25519")

const separator = ":"
const malformedEncoded = "malformed encoded: %s"

type createKpAction struct{}

func (a createKpAction) Execute(ctx node.Context) error {
	kp := key.NewKeyPair(suites.MustFind("Ed25519"))

	sk, err := kp.Private.MarshalBinary()
	if err != nil {
		return xerrors.Errorf("failed to marshal secret key: %v", err)
	}

	pk, err := kp.Public.MarshalBinary()
	if err != nil {
		return xerrors.Errorf("failed to marshal public key: %v", err)
	}

	fmt.Fprintf(ctx.Out, "%v:%v", string(sk), string(pk))

	return nil
}

type revealAction struct{}

func (a revealAction) Execute(ctx node.Context) error {
	encrypted := ctx.Flags.String("encrypted")
	_, cs, err := decodeEncrypted(encrypted)
	if err != nil {
		return xerrors.Errorf("failed to decode encrypted str: %v", err)
	}

	xpk := ctx.Flags.String("xpk")
	xhatenc, err := decodePublicKey(xpk)
	if err != nil {
		return xerrors.Errorf("failed to reencrypt: %v", err)
	}

	dpk := ctx.Flags.String("dpk")
	dpubk, err := decodePublicKey(dpk)
	if err != nil {
		return xerrors.Errorf("failed to decode public key str: %v", err)
	}

	usk := ctx.Flags.String("usk")
	usersk, err := decodePrivateKey(usk)
	if err != nil {
		return xerrors.Errorf("failed to reencrypt: %v", err)
	}

	msg, err := reveal(cs, xhatenc, dpubk, usersk)
	if err != nil {
		fmt.Fprintf(ctx.Out, "couldn't reveal message. %v", err)
		return err
	}
	fmt.Fprint(ctx.Out, string(msg))

	return nil
}

func decodePrivateKey(str string) (sk kyber.Scalar, err error) {
	skbuff, err := hex.DecodeString(str)
	if err != nil {
		return nil, xerrors.Errorf(malformedEncoded, str)
	}

	sk = suite.Scalar()

	err = sk.UnmarshalBinary(skbuff)
	if err != nil {
		return nil, xerrors.Errorf("failed to unmarshal sk: %v", err)
	}

	return sk, nil
}

func decodePublicKey(str string) (pk kyber.Point, err error) {
	pkbuff, err := hex.DecodeString(str)
	if err != nil {
		return nil, xerrors.Errorf(malformedEncoded, str)
	}

	pk = suite.Point()

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

	return k, cs, nil
}
