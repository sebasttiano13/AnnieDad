package models

type (
	// User is a client avatar
	User struct {
		ID             int
		Name           string `json:"name" valid:"required,type(string)"`
		HashedPassword string `json:"password" db:"password" valid:"required,type(string)"`
		RegisteredAT   string
	}
)
