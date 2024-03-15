package controller

import (
	"go.dedis.ch/dela"
	"go.dedis.ch/kyber/v3"
	"golang.org/x/xerrors"
)

// reveal decrypts a reencrypted message.
func reveal(
	XhatEnc kyber.Point,
	dkgPk kyber.Point,
	userPrivateKey kyber.Scalar,
	Cs []kyber.Point,
) ([]byte, error) {
	dela.Logger.Info().Msgf("XhatEnc:%v", XhatEnc)
	dela.Logger.Info().Msgf("dkgPk:%v", dkgPk)
	dela.Logger.Info().Msgf("Cs:%v", Cs)

	xcInv := suite.Scalar().Neg(userPrivateKey)
	XhatDec := suite.Point().Mul(xcInv, dkgPk)
	Xhat := suite.Point().Add(XhatEnc, XhatDec)
	XhatInv := suite.Point().Neg(Xhat)

	msg := make([]byte, 0, 128*len(Cs))

	// Decrypt Cs to keyPointHat
	for _, C := range Cs {
		keyPointHat := suite.Point().Add(C, XhatInv)
		keyPart, err := keyPointHat.Data()
		if err != nil {
			e := xerrors.Errorf("Error while decrypting Cs: %v", err)
			dela.Logger.Error().Err(e).Msg("Failed revealing message")
			return nil, e
		}
		msg = append(msg, keyPart...)
	}

	return msg, nil
}
