package user

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
	"go.dedis.ch/kyber/v3"
)

const smcServer = "http://localhost:40001"

func SmcGetKey() kyber.Point {
	resp, err := http.Get(smcServer + "/smc/pubkey")
	if err != nil {
		log.Fatal().Msgf("error: %v", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	// Decode the response
	var data []byte
	err = decoder.Decode(&data)
	if err != nil {
		log.Fatal().Msgf("error: %v", err)
	}

	// Unmarshal the response
	pk := suite.Point()
	err = pk.UnmarshalBinary(data)
	if err != nil {
		log.Error().Msgf("error decoding response: %v", err)
	}

	return pk
}
