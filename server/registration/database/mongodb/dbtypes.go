package mongodb

// Document is a database struct for the registration service
type Document struct {
	Name       string `bson:"name"`
	Passport   string `bson:"passport"`
	Role       uint   `bson:"role"`
	Picture    []byte `bson:"picture"`
	Hash       []byte `bson:"hash"`
	Registered bool   `bson:"registered"`
}
