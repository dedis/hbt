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

type dbAccess struct {
	client *mongo.Client
}

// NewDBAccess creates a new user access to the DB
func NewDBAccess() (database.Database, error) {
	// Initialize userDb DB access
	credentials := options.Credential{
		AuthSource:    "admin",
		AuthMechanism: "SCRAM-SHA-1",
		Username:      config.AppConfig.UserName,
		Password:      config.AppConfig.UserPassword,
	}
	userOpts := options.Client().ApplyURI(config.AppConfig.MongodbURI).SetAuth(credentials)
	client, err := mongo.Connect(context.TODO(), userOpts)
	if err != nil {
		return dbAccess{nil}, err
	}

	return dbAccess{client}, nil
}

// Create creates a new document in the DB
func (d dbAccess) Create(data *registry.RegistrationData) (*registry.RegistrationID, error) {
	doc := Document{
		Name:       data.Name,
		Passport:   data.Passport,
		Role:       data.Role,
		Picture:    data.Picture,
		Registered: false,
	}

	db := d.client.Database("registry")
	c := db.Collection("documents")
	result, err := c.InsertOne(context.Background(), doc)

	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		id, err := oid.MarshalJSON()
		if err != nil {
			return nil, err
		}

		regID := registry.RegistrationID{
			ID: id,
		}

		return &regID, nil

	}

	return nil, errors.New("could not marshal object id")
}

// Read reads a document from the DB
// it is used to get the registered value of a document
func (d dbAccess) Read(id registry.RegistrationID) (
	*registry.RegistrationData,
	error,
) {
	var doc Document

	err := d.client.Database("registration").Collection("documents").FindOne(context.Background(),
		id).Decode(&doc)

	if err != nil {
		return nil, err
	}

	data := registry.RegistrationData{
		Name:       doc.Name,
		Passport:   doc.Passport,
		Role:       doc.Role,
		Registered: doc.Registered,
	}

	return &data, nil
}

// Update updates a document in the DB
func (d dbAccess) Update(
	id registry.RegistrationID,
	reg *registry.RegistrationData,
) error {
	var doc Document

	err := d.client.Database("registration").Collection("documents").FindOne(context.Background(),
		id).Decode(&doc)
	if err != nil {
		return err
	}

	result, err := d.client.Database("registration").Collection("documents").UpdateOne(context.Background(),
		id, reg)
	if err != nil {
		return err
	}

	if result.ModifiedCount != 1 {
		return errors.New("could not update document")
	}

	return nil
}

// Delete updates a document in the DB
func (d dbAccess) Delete(id registry.RegistrationID) error {
	var doc Document

	err := d.client.Database("registration").Collection("documents").FindOne(context.Background(),
		id).Decode(&doc)
	if err != nil {
		return err
	}

	result, err := d.client.Database("registration").Collection("documents").DeleteOne(context.Background(),
		id)
	if err != nil {
		return err
	}
	if result.DeletedCount != 1 {
		return errors.New("could not delete document")
	}

	return nil
}

// Disconnect disconnects the user from the DB
func (d dbAccess) Disconnect() error {
	err := d.client.Disconnect(context.Background())
	if err != nil {
		panic(err)
	}

	return nil
}
