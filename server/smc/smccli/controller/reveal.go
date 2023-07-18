package controller

import (
	"go.dedis.ch/dela"
	"go.dedis.ch/kyber/v3"
	"golang.org/x/xerrors"
)

// reveal decrypts a reencrypted message.
func reveal(XhatEnc kyber.Point, dkgPk kyber.Point, userPrivateKey kyber.Scalar, Cs []kyber.Point) ([]byte, error) {
	dela.Logger.Info().Msgf("XhatEnc:%v", XhatEnc)
	dela.Logger.Info().Msgf("dkgPk:%v", dkgPk)
	dela.Logger.Info().Msgf("Cs:%v", Cs)

	xcInv := suite.Scalar().Neg(userPrivateKey)
	dela.Logger.Debug().Msgf("xcInv:%v", xcInv)

	sum := suite.Scalar().Add(userPrivateKey, xcInv)
	dela.Logger.Debug().Msgf("xc + xcInv: %v", sum)

	XhatDec := suite.Point().Mul(xcInv, dkgPk)
	dela.Logger.Debug().Msgf("XhatDec:%v", XhatDec)

	Xhat := suite.Point().Add(XhatEnc, XhatDec)
	dela.Logger.Debug().Msgf("Xhat:%v", Xhat)

	XhatInv := suite.Point().Neg(Xhat)
	dela.Logger.Debug().Msgf("XhatInv:%v", XhatInv)

	msg := make([]byte, 0, 128*len(Cs))

	// Decrypt Cs to keyPointHat
	for _, C := range Cs {
		dela.Logger.Debug().Msgf("C:%v", C)

		keyPointHat := suite.Point().Add(C, XhatInv)
		dela.Logger.Debug().Msgf("keyPointHat:%v", keyPointHat)

		keyPart, err := keyPointHat.Data()
		dela.Logger.Debug().Msgf("keyPart:%v", keyPart)

		if err != nil {
			e := xerrors.Errorf("Error while decrypting Cs: %v", err)
			dela.Logger.Error().Msg(e.Error())
			return nil, e
		}
		msg = append(msg, keyPart...)
	}

	return msg, nil
}
