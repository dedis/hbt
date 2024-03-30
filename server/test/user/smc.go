package user

import (
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
	"go.dedis.ch/kyber/v3"
)

// Define a struct to unmarshal the JSON response
type Response struct {
	PubKey []byte `json:"pubkey"`
}

const smcServer = "http://localhost:41001"

func SmcGetKey() kyber.Point {
	resp, err := http.Get(smcServer + "/smc/pubkey")
	if err != nil {
		log.Fatal().Msgf("error: %v", err)
	}
	defer resp.Body.Close()

	// Check if the status code is OK
	if resp.StatusCode != http.StatusOK {
		log.Fatal().Msgf("Error received from server: %v", resp)
	}

	// Read the response body into a byte array
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal().Msgf("error: %v", err)
	}

	// Unmarshal the response
	pk := suite.Point()
	err = pk.UnmarshalBinary(body)
	if err != nil {
		log.Error().Msgf("Error unmarshaling the pubkey: %v", err)
	}

	return pk
}
