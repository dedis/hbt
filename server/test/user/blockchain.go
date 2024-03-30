package user

import (
	"bytes"
	"encoding/hex"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/rs/zerolog/log"
	"go.dedis.ch/dela"
	"go.dedis.ch/hbt/server/registry/registry"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/suites"
	"golang.org/x/xerrors"
)

const blockchainServer = "http://localhost:40001"

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

	var body bytes.Buffer
	w := multipart.NewWriter(&body)

	fw, err := w.CreateFormField("secret")
	if err != nil {
		log.Fatal().Msgf("error: %v", err)
	}
	_, err = io.Copy(fw, bytes.NewReader([]byte(encryptedSecret)))
	if err != nil {
		log.Fatal().Msgf("error: %v", err)
	}

	fw, err = w.CreateFormField("id")
	if err != nil {
		log.Fatal().Msgf("error: %v", err)
	}
	_, err = io.Copy(fw, bytes.NewReader(id.ID))
	if err != nil {
		log.Fatal().Msgf("error: %v", err)
	}

	w.Close()

	req, err := http.NewRequest(http.MethodPost, blockchainServer+"/secret", &body)
	if err != nil {
		log.Fatal().Msgf("error: %v", err)
	}

	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
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
