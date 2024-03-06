package user

import (
	"encoding/hex"
	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"
	"go.dedis.ch/dela"
	"go.dedis.ch/hbt/server/registration/registry"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/suites"
	"golang.org/x/xerrors"
)

const blockchainServer = "localhost:3003"

// suite is the Kyber suite for Pedersen.
var suite = suites.MustFind("Ed25519")

func BlockchainEncryptAndAddSecret(
	key kyber.Point,
	secret []byte,
	id registry.RegistrationID,
) string {
	// Encrypt the secret
	encryptedSecret, err := encryptToKCs(key, secret)
	if err != nil {
		log.Fatal().Msgf("error: %v", err)
	}

	// Add the secret to the blockchain
	resp, err := http.PostForm(blockchainServer+"/secret",
		url.Values{
			"secret": {encryptedSecret},
			"id":     {string(id.ID)},
		})

	if err != nil {
		log.Fatal().Msgf("error: %v", err)
	}

	defer resp.Body.Close()

	return encryptedSecret
}

func encryptToKCs(key kyber.Point, msg []byte) (string, error) {
	// ElGamal-encrypt the point to produce ciphertext (K,C).
	r := suite.Scalar().Pick(suite.RandomStream())

	K := suite.Point().Mul(r, nil)
	dela.Logger.Debug().Msgf("K: %v", K.String())

	C := suite.Point().Mul(r, key)
	dela.Logger.Debug().Msgf("C: %v", C)

	// S: ephemeral DH shared secret
	S := suite.Point().Mul(r, key)
	dela.Logger.Debug().Msgf("S: %v", S.String())

	Cs := make([]kyber.Point, 0, 16)
	for len(msg) > 0 {
		kp := suite.Point().Embed(msg, suite.RandomStream())
		dela.Logger.Debug().Msgf("kp: %v", kp.String())

		// message blinded with secret
		c := suite.Point().Add(C, kp)
		dela.Logger.Debug().Msgf("c: %v", c)

		Cs = append(Cs, c)
		dela.Logger.Debug().Msgf("Cs: %v", Cs)

		msg = msg[min(len(msg), kp.EmbedLen()):]
	}

	return encodeEncrypted(K, Cs)
}

// from pedersen/controller/action
const separator = ":"

func encodeEncrypted(k kyber.Point, cs []kyber.Point) (string, error) {
	kbuff, err := k.MarshalBinary()
	if err != nil {
		return "", xerrors.Errorf("failed to marshal k: %v", err)
	}

	encoded := hex.EncodeToString(kbuff)

	for _, c := range cs {
		cbuff, err := c.MarshalBinary()
		if err != nil {
			return "", xerrors.Errorf("failed to marshal c: %v", err)
		}
		encoded += separator + hex.EncodeToString(cbuff)
	}

	return encoded, nil
}
