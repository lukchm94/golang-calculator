package appEvent

type AppEventDetailType string

const (
	Calculation AppEventDetailType = "calculation"
	Login       AppEventDetailType = "login"
)

func (d AppEventDetailType) String() string {
	return string(d)
}

func (d AppEventDetailType) IsValid() bool {
	switch d {
	case Calculation, Login:
		return true
	}
	return false
}
