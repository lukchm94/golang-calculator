package userDomain

type Role string

const (
	Admin Role = "admin"
	Guest Role = "guest"
)

func (r Role) IsValid() bool {
	switch r {
	case Admin, Guest:
		return true
	default:
		return false
	}
}

func (r Role) String() string {
	return string(r)
}
