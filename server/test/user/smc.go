package user

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
	"go.dedis.ch/kyber/v3"
)

func SmcGetKey() kyber.Point {
	resp, err := http.Get("localhost:3000/key")
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
