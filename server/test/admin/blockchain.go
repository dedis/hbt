package admin

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
	"go.dedis.ch/hbt/server/registry/registry"
	"go.dedis.ch/hbt/server/smc"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/suites"
)

const blockchainServer = "http://localhost:40001"

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

	// Reading the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Msgf("Error reading response body: %v", err)
		return nil
	}

	var items []registry.RegistrationID
	// Checking if the request was successful (status code 200)
	if resp.StatusCode == http.StatusOK {
		// Parsing JSON data
		if err := json.Unmarshal(body, &items); err != nil {
			log.Error().Msgf("Error parsing JSON: %v", err)
			return nil
		}

		// Printing the list of IDs
		log.Info().Msg("List of IDs:")
		for i, item := range items {
			log.Info().Msgf("ID[%v] = %v", i, item.ID)
		}
	} else {
		log.Error().Msgf("Failed to fetch items. Status code:%v", resp.StatusCode)
	}

	return items
}

// BlockchainGetDocument polls the blockchain to get the encrypted document
func BlockchainGetSecret(id registry.RegistrationID, pk kyber.Point) (smc.Secret, []byte) {
	encodedPk, err := pk.MarshalBinary()
	if err != nil {
		log.Fatal().Msgf("error: %v", err)
	}

	resp, err := http.Get(blockchainServer + "/secret?pubkey=" + string(encodedPk) + "&id=" + string(id.ID))
	if err != nil {
		log.Fatal().Msgf("error: %v", err)
	}

	defer resp.Body.Close()

	// Decode the response
	var secret smc.Secret
	err = json.NewDecoder(resp.Body).Decode(&secret)
	if err != nil {
		log.Error().Msgf("error decoding response: %v", err)
	}

	return secret, nil
}
