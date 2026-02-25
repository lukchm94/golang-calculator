package userDomain

type UserRepository interface {
	GetUserByID(id string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	CreateUser(user *User) error
	GenerateUserID() string
}
