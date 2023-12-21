package registry

type DocId []byte

// RegistrationData contains the data for a registration
type RegistrationData struct {
	Name       string `json:"name"`
	Passport   string `json:"passport"`
	Role       uint   `json:"role"`
	Picture    []byte `json:"picture"`
	Registered bool   `json:"registered"`
}

// RegistrationId contains the reference to the document in the database
type RegistrationId struct {
	Id DocId `json:"doc_id"`
}

// Document is a database struct for the registration service
type Document struct {
	Name       string `bson:"name"`
	Passport   string `bson:"passport"`
	Role       uint   `bson:"role"`
	Picture    []byte `bson:"picture"`
	Registered bool   `bson:"registered"`
}
