package controller

import (
	"go.dedis.ch/dela"
	"go.dedis.ch/kyber/v3"
	"golang.org/x/xerrors"
)

// reveal decrypts a reencrypted message.
func reveal(Cs []kyber.Point, XhatEnc kyber.Point, dkgPk kyber.Point, Sk kyber.Scalar) ([]byte, error) {
	dela.Logger.Debug().Msgf("DKG pubK:%v", dkgPk)
	dela.Logger.Debug().Msgf("XhatEnc:%v", XhatEnc)
	dela.Logger.Debug().Msgf("xc:%v", Sk)

	xcInv := suite.Scalar().Neg(Sk)
	dela.Logger.Debug().Msgf("xcInv:%v", xcInv)

	sum := suite.Scalar().Add(Sk, xcInv)
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
