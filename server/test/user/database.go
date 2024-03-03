package user

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"
	"go.dedis.ch/hbt/server/registration/registry"
	"go.dedis.ch/hbt/server/test/key"
)

const registrationServer = "localhost:3000"

// RegistrationAdd adds a new registration to the registry
func RegistrationAdd(data registry.RegistrationData, symKey []byte) registry.RegistrationID {
	// Encrypt the data
	encrypted, err := encryptRegistrationData(data, symKey)
	if err != nil {
		log.Fatal().Msgf("error: %v", err)
	}

	// Add the encrypted document to the registry
	resp, err := http.PostForm(registrationServer+"/document",
		url.Values{
			"name":       {string(encrypted.Name)},
			"passport":   {string(encrypted.Passport)},
			"role":       {string(encrypted.Role)},
			"picture":    {string(encrypted.Picture)},
			"registered": {string(encrypted.Registered)},
		})
	if err != nil {
		log.Fatal().Msgf("error: %v", err)
	}

	defer resp.Body.Close()

	// Decode the response
	var docid registry.RegistrationID
	err = json.NewDecoder(resp.Body).Decode(&docid)
	if err != nil {
		log.Error().Msgf("error decoding response: %v", err)
	}

	return docid
}

// RegistrationGet polls the data to see if registered
func RegistrationGet(docid registry.RegistrationID, symKey []byte) registry.RegistrationData {
	resp, err := http.Get(registrationServer + "/document?id=" + string(docid.ID))
	if err != nil {
		log.Fatal().Msgf("error: %v", err)
	}

	defer resp.Body.Close()

	// Decode the response
	var encrypted registry.EncryptedData
	err = json.NewDecoder(resp.Body).Decode(&encrypted)
	if err != nil {
		log.Error().Msgf("error decoding response: %v", err)
	}

	// Decrypt the data
	data, err := decryptRegistrationData(encrypted, symKey)
	if err != nil {
		log.Fatal().Msgf("error: %v", err)
	}

	return data
}

// RegistrationDelete deletes the registration data from the database
func RegistrationDelete(docid registry.RegistrationID) error {
	req, err := http.NewRequest("DELETE", registrationServer+"/document?id="+string(docid.ID), nil)
	if err != nil {
		log.Fatal().Msgf("error: %v", err)
	}

	ctx := context.Background()
	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal().Msgf("error: %v", err)
	}

	defer resp.Body.Close()

	return err
}

// ---------------------------------------------------------------------------
// The following functions are used to encrypt and decrypt the registration

func encryptRegistrationData(data registry.RegistrationData, symKey []byte) (
	registry.EncryptedData,
	error,
) {
	// Convert the struct to a byte array
	buf := new(bytes.Buffer)

	// encrypt the data.Name
	err := binary.Write(buf, binary.LittleEndian, data.Name)
	if err != nil {
		return registry.EncryptedData{}, err
	}

	encName, err := key.Encrypt(symKey, buf.Bytes())
	if err != nil {
		return registry.EncryptedData{}, err
	}

	// encrypt the data.Passport
	err = binary.Write(buf, binary.LittleEndian, data.Name)
	if err != nil {
		return registry.EncryptedData{}, err
	}

	encPassport, err := key.Encrypt(symKey, buf.Bytes())
	if err != nil {
		return registry.EncryptedData{}, err
	}

	// encrypt the data.Role
	err = binary.Write(buf, binary.LittleEndian, data.Role)
	if err != nil {
		return registry.EncryptedData{}, err
	}

	encRole, err := key.Encrypt(symKey, buf.Bytes())
	if err != nil {
		return registry.EncryptedData{}, err
	}

	// encrypt the data.Picture
	err = binary.Write(buf, binary.LittleEndian, data.Picture)
	if err != nil {
		return registry.EncryptedData{}, err
	}

	encPicture, err := key.Encrypt(symKey, buf.Bytes())
	if err != nil {
		return registry.EncryptedData{}, err
	}

	// encrypt the data.Registered
	err = binary.Write(buf, binary.LittleEndian, data.Registered)
	if err != nil {
		return registry.EncryptedData{}, err
	}

	encRegistered, err := key.Encrypt(symKey, buf.Bytes())
	if err != nil {
		return registry.EncryptedData{}, err
	}

	return registry.EncryptedData{
		Name:       encName,
		Passport:   encPassport,
		Picture:    encPicture,
		Role:       encRole,
		Registered: encRegistered,
	}, nil
}

func decryptRegistrationData(encrypted registry.EncryptedData, symKey []byte) (
	registry.RegistrationData,
	error,
) {
	// Decrypt the data.Name
	decName, err := key.Decrypt(symKey, encrypted.Name)
	if err != nil {
		return registry.RegistrationData{}, err

	}

	// Decrypt the data.Passport
	decPassport, err := key.Decrypt(symKey, encrypted.Passport)
	if err != nil {
		return registry.RegistrationData{}, err
	}

	// Decrypt the data.Role
	decRole, err := key.Decrypt(symKey, encrypted.Role)
	if err != nil {
		return registry.RegistrationData{}, err
	}

	// Decrypt the data.Picture
	decPicture, err := key.Decrypt(symKey, encrypted.Picture)
	if err != nil {
		return registry.RegistrationData{}, err
	}

	// Decrypt the data.Registered
	decRegistered, err := key.Decrypt(symKey, encrypted.Registered)
	if err != nil {
		return registry.RegistrationData{}, err
	}

	return registry.RegistrationData{
		Name:       string(decName),
		Passport:   string(decPassport),
		Picture:    decPicture,
		Role:       binary.LittleEndian.Uint64(decRole),
		Registered: binary.LittleEndian.Uint64(decRegistered) == 1,
	}, nil

}
