package user

import (
	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"
	"go.dedis.ch/hbt/server/registration/registry"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/suites"
)

const blockchainServer = "localhost:4000"

// suite is the Kyber suite for Pedersen.
var suite = suites.MustFind("Ed25519")

func BlockchainEncryptAndAddSecret(key kyber.Point, secret []byte, id registry.RegistrationID) {
	// Encrypt the secret
	encryptedSecret := suite.Point().Mul(suite.Scalar().SetBytes(secret), key)

	// Add the secret to the blockchain
	resp, err := http.PostForm(blockchainServer+"/secret",
		url.Values{
			"secret": {encryptedSecret.String()},
			"id":     {string(id.ID)},
		})

	if err != nil {
		log.Fatal().Msgf("error: %v", err)
	}

	defer resp.Body.Close()
}
