package mongodb

import (
	"context"
	"errors"

	"go.dedis.ch/hbt/server/registration/config"
	"go.dedis.ch/hbt/server/registration/database"
	"go.dedis.ch/hbt/server/registration/registry"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserDbAccess struct {
	client *mongo.Client
}

// NewUserDbAccess creates a new user access to the DB
func NewUserDbAccess() database.Database {
	// Initialize userDb DB access
	userCredentials := options.Credential{
		Username: config.AppConfig.UserName,
		Password: config.AppConfig.UserPassword,
	}
	userOpts := options.Client().ApplyURI(config.AppConfig.MongoDbUri).SetAuth(userCredentials)
	userDb, err := mongo.Connect(context.TODO(), userOpts)
	if err != nil {

		return nil
	}

	return UserDbAccess{
		client: userDb,
	}
}

// Create creates a new document in the DB
func (u UserDbAccess) Create(doc registry.RegistrationData) (registry.DocId, error) {
	result, err := u.client.Database("registration").Collection("documents").InsertOne(context.Background(),
		doc)

	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return oid.MarshalJSON()
	}

	return nil, errors.New("could not marshal object id")
}

// Read reads a document from the DB
func (u UserDbAccess) Read(docId registry.DocId) (*registry.RegistrationData, error) {
	return nil, nil
}

// Update updates a document in the DB
func (u UserDbAccess) Update(docId registry.DocId, reg *registry.RegistrationData) error {
	return nil
}

// Delete updates a document in the DB
func (u UserDbAccess) Delete(docId registry.DocId) error {
	return nil
}

// Disconnect disconnects the user from the DB
func (u UserDbAccess) Disconnect() error {
	err := u.client.Disconnect(context.Background())
	if err != nil {
		panic(err)
	}

	return nil
}
