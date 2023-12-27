package mongodb

import (
	"bytes"
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
func (u UserDbAccess) Create(data registry.RegistrationData) (*registry.RegistrationId, error) {
	doc := Document{
		Name:       data.Name,
		Passport:   data.Passport,
		Role:       data.Role,
		Picture:    data.Picture,
		Hash:       data.Hash,
		Registered: false,
	}

	result, err := u.client.Database("registration").Collection("documents").InsertOne(context.Background(),
		doc)

	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		id, err := oid.MarshalJSON()
		if err != nil {
			return nil, err
		}

		registrationId := registry.RegistrationId{
			Id: id,
		}

		return &registrationId, nil

	}

	return nil, errors.New("could not marshal object id")
}

// Read reads a document from the DB
// it is used to get the registered value of a document
func (u UserDbAccess) Read(docId registry.RegistrationId, hash []byte) (
	*registry.RegistrationData,
	error,
) {
	var doc Document

	err := u.client.Database("registration").Collection("documents").FindOne(context.Background(),
		docId).Decode(&doc)

	if err != nil {
		return nil, err
	}

	if hash != nil {
		if !bytes.Equal(hash, doc.Hash) {
			return nil, errors.New("hashes do not match")
		}
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
func (u UserDbAccess) Update(
	docId registry.RegistrationId,
	hash []byte,
	reg *registry.RegistrationData,
) error {
	var doc Document

	err := u.client.Database("registration").Collection("documents").FindOne(context.Background(),
		docId).Decode(&doc)
	if err != nil {
		return err
	}

	if !bytes.Equal(hash, doc.Hash) {
		return errors.New("hashes do not match")
	}

	reg.Name = doc.Name
	reg.Passport = doc.Passport
	reg.Role = doc.Role
	reg.Picture = doc.Picture
	reg.Registered = false

	result, err := u.client.Database("registration").Collection("documents").UpdateOne(context.Background(),
		docId, reg)
	if err != nil {
		return err
	}

	if result.ModifiedCount != 1 {
		return errors.New("could not update document")
	}

	return nil
}

// Delete updates a document in the DB
func (u UserDbAccess) Delete(docId registry.RegistrationId, hash []byte) error {
	var doc Document

	err := u.client.Database("registration").Collection("documents").FindOne(context.Background(),
		docId).Decode(&doc)
	if err != nil {
		return err
	}

	if !bytes.Equal(hash, doc.Hash) {
		return errors.New("hashes do not match")
	}

	result, err := u.client.Database("registration").Collection("documents").DeleteOne(context.Background(),
		docId)
	if err != nil {
		return err
	}
	if result.DeletedCount != 1 {
		return errors.New("could not delete document")
	}

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
