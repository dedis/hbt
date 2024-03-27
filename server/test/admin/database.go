package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
	"go.dedis.ch/hbt/server/registry/registry"
)

const registrationServer = "localhost:3001"

func RegistrationAdminGetDocument(docid registry.RegistrationID) registry.RegistrationData {
	resp, err := http.Get(registrationServer + "/admin/document?id=" + string(docid.ID))
	if err != nil {
		log.Fatal().Msgf("error: %v", err)
	}

	defer resp.Body.Close()

	// Decode the response
	var data registry.RegistrationData
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Error().Msgf("error decoding response: %v", err)
	}

	return data
}

func RegistrationAdminUpdateDocument(docid registry.RegistrationID) error {
	resp, err := http.Get(registrationServer + "/admin/document?id=" + string(docid.ID))
	if err != nil {
		log.Fatal().Msgf("error: %v", err)
	}

	defer resp.Body.Close()

	// Decode the response
	var data registry.RegistrationData
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Error().Msgf("error decoding response: %v", err)
	}

	data.Registered = true
	out, err := json.Marshal(data)
	if err != nil {
		log.Fatal().Msgf("error: %v", err)
	}

	req, err := http.NewRequest(http.MethodPut,
		"localhost:3000/admin/document?id="+string(docid.ID),
		bytes.NewBuffer(out))
	if err != nil {
		log.Fatal().Msgf("error: %v", err)
	}

	ctx := context.Background()
	req = req.WithContext(ctx)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Error().Msgf("response: %v", resp)
	}

	defer resp.Body.Close()

	return err
}

func RegistrationAdminDeleteDocument(docid registry.RegistrationID) error {
	req, err := http.NewRequest(http.MethodDelete,
		"localhost:3000/admin/document?id="+string(docid.ID), nil)
	if err != nil {
		log.Fatal().Msgf("error: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal().Msgf("error: %v", err)
	}

	defer resp.Body.Close()

	return err
}
