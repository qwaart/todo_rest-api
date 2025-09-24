package auth

type APIKey struct {
	ID		int64
	Key 	string // hashed key
	Owner 	string
	Revoked	bool
}

type Permission struct {
	ID int64
	Name string
}