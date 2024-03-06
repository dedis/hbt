package mongodb

// Document is a database struct for the registration service
type Document struct {
	Name       string `bson:"name"`
	Passport   string `bson:"passport"`
	Role       uint64 `bson:"role"`
	Picture    []byte `bson:"picture"`
	Registered bool   `bson:"registered"`
}
