package admin

import (
	"net/http"

	"github.com/rs/zerolog/log"
	"go.dedis.ch/hbt/server/registration/registry"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/suites"
)

const blockchainServer = "localhost:4000"

// suite is the Kyber suite for Pedersen.
var suite = suites.MustFind("Ed25519")

// BlockchainGetDocs polls the blockchain to get the list of encrypted documents
// adminPubkey is the public key of the admin and is used for audit purpose
func BlockchainGetDocIDs(adminPubkey kyber.Point) []registry.RegistrationID {
	encoded, err := adminPubkey.MarshalBinary()
	if err != nil {
		log.Fatal().Msgf("error: %v", err)
	}

	resp, err := http.Get(blockchainServer + "/secret/list?pubkey=" + string(encoded))
	if err != nil {
		log.Fatal().Msgf("error: %v", err)
	}

	defer resp.Body.Close()

	// Decode the response
	var data []string

	// TODO: Decode the response and return the list of doc IDs
}
