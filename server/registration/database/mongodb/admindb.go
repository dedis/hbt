package mongodb

import (
	"context"
	"errors"

	"go.dedis.ch/hbt/server/registration/config"
	"go.dedis.ch/hbt/server/registration/registry"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AdminDbAccess struct {
	client *mongo.Client
}

// NewAdminDbAccess creates a new admin access to the DB
func NewAdminDbAccess() *AdminDbAccess {
	// Initialize adminDb DB access
	adminCredentials := options.Credential{
		Username: config.AppConfig.UserName,
		Password: config.AppConfig.UserPassword,
	}
	adminOpts := options.Client().ApplyURI(config.AppConfig.MongoDbUri).SetAuth(adminCredentials)
	adminDb, err := mongo.Connect(context.TODO(), adminOpts)
	if err != nil {

		return nil
	}

	return &AdminDbAccess{
		client: adminDb,
	}
}

// Create creates a new document in the DB
func (a AdminDbAccess) Create(doc registry.RegistrationData) (*registry.RegistrationId, error) {
	return nil, errors.New("admin cannot create user documents")
}

// Read reads a document from the DB
func (a AdminDbAccess) Read(docId registry.RegistrationId, hash []byte) (
	*registry.RegistrationData,
	error,
) {
	return nil, nil
}

// Update updates a document in the DB
func (a AdminDbAccess) Update(
	docId registry.RegistrationId,
	hash []byte,
	reg *registry.RegistrationData,
) error {
	return nil
}

// Delete updates a document in the DB
func (a AdminDbAccess) Delete(docId registry.RegistrationId, hash []byte) error {
	return nil
}

// Disconnect disconnects the admin from the DB
func (a *AdminDbAccess) Disconnect() error {
	err := a.client.Disconnect(context.Background())
	if err != nil {
		panic(err)
	}

	return nil
}
