package database

import "go.dedis.ch/hbt/server/registration/registry"

// Database defines a generic CRUD interface to the database
type Database interface {
	// Create creates a new document in the database
	Create(registry.RegistrationData) (*registry.RegistrationId, error)

	// Read retrieves a document from the database
	// it takes the document ID as an argument
	// and returns the document
	Read(registry.RegistrationId, []byte) (*registry.RegistrationData, error)

	// Update updates a document in the database
	// it takes the document ID and the updated document as an argument
	Update(registry.RegistrationId, []byte, *registry.RegistrationData) error

	// DeleteDocument deletes a document from the database
	// it takes the document ID as an argument
	Delete(registry.RegistrationId, []byte) error

	// Disconnect disconnects from the database
	Disconnect() error
}
