package admin

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
	"go.dedis.ch/dela"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/suites"
	"golang.org/x/xerrors"
)

const smcServer = "localhost:3002"

// SmcReencryptSecret re-encrypts the secret with the new public key
// and returns a xhatenc value that can be used to reveal the secret
// first argument is supposed to be the proof
func SmcReencryptSecret(_ []byte, pk kyber.Point, secret string) (kyber.Point, error) {
	resp, err := http.Post(smcServer+"/reencrypt", "application/json",
		bytes.NewBuffer([]byte(`{"pubk": "`+encodePublickey(pk)+`", "encrypted": "`+secret+`"}`)))
	if err != nil {
		log.Error().Msgf("error: %v", err)
		return nil, err
	}

	defer resp.Body.Close()

	// Decode the response
	var xhatenc string
	err = json.NewDecoder(resp.Body).Decode(&xhatenc)
	if err != nil {
		log.Error().Msgf("error decoding response: %v", err)
		return nil, err
	}

	xhatencbuff, err := decodeReencrypted(xhatenc)
	if err != nil {
		log.Error().Msgf("error decoding response: %v", err)
		return nil, err
	}

	return xhatencbuff, nil
}

func encodePublickey(pk kyber.Point) string {
	pkbuff, err := pk.MarshalBinary()
	if err != nil {
		return ""
	}

	return hex.EncodeToString(pkbuff)
}

// func encodeSecret(secret []byte) string {
// 	return hex.EncodeToString(secret)
// }

func decodeReencrypted(xhatencstring string) (kyber.Point, error) {
	xhatencbuff, err := hex.DecodeString(xhatencstring)
	if err != nil {
		return nil, err
	}

	xhatenc := suite.Point()
	err = xhatenc.UnmarshalBinary(xhatencbuff)
	if err != nil {
		return nil, err
	}

	return xhatenc, nil
}

func SmcGetKey() kyber.Point {
	resp, err := http.Get(smcServer + "/key")
	if err != nil {
		log.Fatal().Msgf("error: %v", err)
	}

	defer resp.Body.Close()

	// Decode the response
	var data kyber.Point

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Error().Msgf("error decoding response: %v", err)
	}

	return data
}

// SmcReveal decrypts a reencrypted message.
func SmcReveal(
	XhatEnc kyber.Point,
	dkgPk kyber.Point,
	userPrivateKey kyber.Scalar,
	secret string,
) ([]byte, error) {
	_, Cs, err := decodeEncrypted(secret)
	if err != nil {
		log.Fatal().Msgf("error: %v", err)
	}

	log.Info().Msgf("XhatEnc:%v", XhatEnc)
	log.Info().Msgf("dkgPk:%v", dkgPk)
	log.Info().Msgf("Cs:%v", Cs)

	suite := suites.MustFind("Ed25519")

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

// straight from pedersen/controller/action
func decodeEncrypted(str string) (kyber.Point, []kyber.Point, error) {
	const separator = ":"
	parts := strings.Split(str, separator)
	if len(parts) < 2 {
		return nil, nil, xerrors.Errorf("malformed encoded: %s", str)
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
