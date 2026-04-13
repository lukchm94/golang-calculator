package userEvents

import (
	userDomain "app/internal/domain/user"
	"encoding/json"
	"time"
)

type LoginEvent struct {
	UserID    string          `json:"userId"`
	Timestamp time.Time       `json:"timestamp"`
	Role      userDomain.Role `json:"role"`
}

func (e LoginEvent) JSON() string {
	data, err := json.Marshal(e)
	if err != nil {
		return "{}"
	}

	return string(data)
}
