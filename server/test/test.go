package main

import (
	"os"

	"github.com/rs/zerolog/log"
	"go.dedis.ch/hbt/server/registration/registry"
	"go.dedis.ch/hbt/server/test/admin"
	"go.dedis.ch/hbt/server/test/key"
	"go.dedis.ch/hbt/server/test/user"
)

const keySize = 32

func main() {
	// create a secret symmetric key
	symKey := key.NewSymetric(keySize)

	// PRETEND TO BE A USER
	// ---------------------------------------------------------

	// create a document and save it encrypted into the database
	doc := createDocument("John Doe", "12AB456789", 0, "test/passport.jpg")
	log.Info().Msg("SUCCESS! created new document")

	// add the document to the registry
	docid := user.RegistrationAdd(doc, symKey)
	log.Info().Msgf("SUCCESS! added document id: %v", docid)

	// get the SMC pub key
	smcKey := user.SmcGetKey()
	log.Info().Msgf("SUCCESS! added document id: %v", docid)

	// add secret = symKey to the blockchain
	user.BlockchainEncryptAndAddSecret(smcKey, symKey, docid)

	// PRETEND TO BE AN ADMIN
	// ---------------------------------------------------------
	// create a new admin asymmetric key pair
	pk, sk := key.NewAsymmetric()

	// fetch the list of docs from the blockchain
	docIDs := admin.BlockchainGetDocIDs(pk)

	for _, id := range docIDs {
		doc := admin.BlockchainGetDocument(id)
		log.Info().Msgf("document: %v", doc)

		reencrypted := admin.SmcReencryptSecret(pk, id)

		encryptedDoc = admin.registrationGetDocument(id)
	}
}

// ---------------------------------------------------------
// helper functions

// create a document from a picture file
func createDocument(name, passport string, role uint64, picture string) registry.RegistrationData {
	// load picture from file named picture
	picData, err := os.ReadFile(picture)
	if err != nil {
		log.Fatal().Msgf("error while reading picture file: %v", err)
	}

	return registry.RegistrationData{
		Name:       name,
		Passport:   passport,
		Role:       role,
		Picture:    picData,
		Registered: false,
	}
}
