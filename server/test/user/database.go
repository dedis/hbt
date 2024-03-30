package user

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
	"go.dedis.ch/hbt/server/registry/registry"
	"go.dedis.ch/hbt/server/test/key"
)

const registrationServer = "http://localhost:3000"

// RegistrationAdd adds a new registration to the registry
func RegistrationAdd(data registry.RegistrationData, symKey []byte) (
	registry.RegistrationID,
	error,
) {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)

	fw, err := w.CreateFormField("pubkey")
	if err != nil {
		return registry.RegistrationID{}, err
	}
	if _, err = io.Copy(fw, bytes.NewReader(symKey)); err != nil {
		return registry.RegistrationID{}, err
	}

	fw, err = w.CreateFormField("name")
	if err != nil {
		return registry.RegistrationID{}, err
	}
	if _, err = io.Copy(fw, bytes.NewReader([]byte(data.Name))); err != nil {
		return registry.RegistrationID{}, err
	}

	fw, err = w.CreateFormField("passport")
	if err != nil {
		return registry.RegistrationID{}, err
	}
	if _, err = io.Copy(fw, bytes.NewReader([]byte(data.Passport))); err != nil {
		return registry.RegistrationID{}, err
	}

	fw, err = w.CreateFormField("role")
	if err != nil {
		return registry.RegistrationID{}, err
	}
	if _, err = io.Copy(fw,
		bytes.NewReader([]byte(strconv.FormatUint(data.Role, 10)))); err != nil {
		return registry.RegistrationID{}, err
	}

	fw, err = w.CreateFormField("registered")
	if err != nil {
		return registry.RegistrationID{}, err
	}
	if _, err = io.Copy(fw,
		bytes.NewReader([]byte(strconv.FormatBool(data.Registered)))); err != nil {
		return registry.RegistrationID{}, err
	}

	fw, err = w.CreateFormFile("portrait", "portrait.jpg")
	if err != nil {
		return registry.RegistrationID{}, err
	}
	if _, err = io.Copy(fw, bytes.NewReader(data.Picture)); err != nil {
		return registry.RegistrationID{}, err
	}

	w.Close()

	req, err := http.NewRequest(http.MethodPost, registrationServer+"/document", &body)
	if err != nil {
		return registry.RegistrationID{}, err
	}

	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return registry.RegistrationID{}, err
	}

	defer resp.Body.Close()

	// Decode the response
	var docid registry.RegistrationID
	err = json.NewDecoder(resp.Body).Decode(&docid)
	if err != nil {
		return registry.RegistrationID{}, err
	}

	return docid, err
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
	req, err := http.NewRequest(http.MethodDelete,
		registrationServer+"/document?id="+string(docid.ID), nil)
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

// func encryptRegistrationData(data registry.RegistrationData, symKey []byte) (
// 	registry.EncryptedData,
// 	error,
// ) {
// 	// Convert the struct to a byte array
// 	buf := new(bytes.Buffer)
//
// 	// encrypt the data.Name
// 	err := binary.Write(buf, binary.LittleEndian, data.Name)
// 	if err != nil {
// 		return registry.EncryptedData{}, err
// 	}
//
// 	encName, err := key.Encrypt(symKey, buf.Bytes())
// 	if err != nil {
// 		return registry.EncryptedData{}, err
// 	}
//
// 	// encrypt the data.Passport
// 	err = binary.Write(buf, binary.LittleEndian, data.Name)
// 	if err != nil {
// 		return registry.EncryptedData{}, err
// 	}
//
// 	encPassport, err := key.Encrypt(symKey, buf.Bytes())
// 	if err != nil {
// 		return registry.EncryptedData{}, err
// 	}
//
// 	// encrypt the data.Role
// 	err = binary.Write(buf, binary.LittleEndian, data.Role)
// 	if err != nil {
// 		return registry.EncryptedData{}, err
// 	}
//
// 	encRole, err := key.Encrypt(symKey, buf.Bytes())
// 	if err != nil {
// 		return registry.EncryptedData{}, err
// 	}
//
// 	// encrypt the data.Picture
// 	err = binary.Write(buf, binary.LittleEndian, data.Picture)
// 	if err != nil {
// 		return registry.EncryptedData{}, err
// 	}
//
// 	encPicture, err := key.Encrypt(symKey, buf.Bytes())
// 	if err != nil {
// 		return registry.EncryptedData{}, err
// 	}
//
// 	// encrypt the data.Registered
// 	err = binary.Write(buf, binary.LittleEndian, data.Registered)
// 	if err != nil {
// 		return registry.EncryptedData{}, err
// 	}
//
// 	encRegistered, err := key.Encrypt(symKey, buf.Bytes())
// 	if err != nil {
// 		return registry.EncryptedData{}, err
// 	}
//
// 	return registry.EncryptedData{
// 		Name:       encName,
// 		Passport:   encPassport,
// 		Picture:    encPicture,
// 		Role:       encRole,
// 		Registered: encRegistered,
// 	}, nil
// }

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
