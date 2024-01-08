package database

import "go.dedis.ch/hbt/server/registration/registry"

// Database defines a generic CRUD interface to the database
type Database interface {
	// Create creates a new document in the database
	// it takes the document as an argument
	// and returns the document ID or an error
	Create(*registry.RegistrationData) (*registry.RegistrationId, error)

	// Read retrieves a document from the database
	// it takes the document ID and the hash as arguments
	// and returns the document or an error
	Read(registry.RegistrationId, []byte) (*registry.RegistrationData, error)

	// Update updates a document in the database
	// it takes the document ID and the updated document as an argument
	// and returns nil or an error
	Update(registry.RegistrationId, *registry.RegistrationData) error

	// Delete deletes a document from the database
	// it takes the document ID and the hash as arguments
	// and returns nil or an error
	Delete(registry.RegistrationId, []byte) error

	// Disconnect disconnects from the database
	Disconnect() error
}
