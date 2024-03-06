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
	doc := createDocument("John Doe", "12AB456789", 0, "passport.jpg")
	log.Info().Msg("SUCCESS! created new document")

	// add the document to the registry
	docid := user.RegistrationAdd(doc, symKey)
	log.Info().Msgf("SUCCESS! added document ID=%v", docid)

	// get the SMC pub key
	smcKey := user.SmcGetKey()
	log.Info().Msgf("SUCCESS! got SMC key: %v", smcKey)

	// add secret = symKey to the blockchain
	secret := user.BlockchainEncryptAndAddSecret(smcKey, symKey, docid)
	log.Info().Msgf("SUCCESS! added secret=%v with ID=%v to blockchain", secret, docid)

	// PRETEND TO BE AN ADMIN
	// ---------------------------------------------------------
	// create a new admin asymmetric key pair
	pk, sk := key.NewAsymmetric()

	// fetch the list of docs from the blockchain
	// give it the admin pub key for audit purpose
	docIDs := admin.BlockchainGetDocIDs(pk)

	for _, id := range docIDs {
		secret := admin.BlockchainGetSecret(id, pk)
		log.Info().Msgf("secret: %v", secret)

		xhatenc, err := admin.SmcReencryptSecret(pk, secret.Data)
		if err != nil {
			log.Fatal().Msgf("error: %v", err)
		}

		smcKeyAdmin := admin.SmcGetKey()
		if smcKey != smcKeyAdmin {
			log.Fatal().Msg("SMC key mismatch")
		}

		// secret.Data = K:Cs in a string format
		symKey2, err := admin.SmcReveal(xhatenc, smcKey, sk, secret.Data)

		if false == compare2ByteArrays(symKey, symKey2) {
			log.Fatal().Msg("symmetric key mismatch")
		}

		// TODO: get the encrypted document from the registry
		// TODO: decrypt the document - optional
		// TODO: update the document status to registered
		// TODO: encrypt the document - optional
		// TODO: save the document back to the registry
	}

	// PRETEND TO BE A USER
	// ---------------------------------------------------------
	// get the document from the registry to see the updated status
	doc2 := user.RegistrationGet(docid, symKey)
	log.Info().Msgf("SUCCESS! got document: %v", doc2)
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

func compare2ByteArrays(a []byte, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
