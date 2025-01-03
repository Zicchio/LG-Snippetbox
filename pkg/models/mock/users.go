package mock

import (
	"time"

	"github.com/Zicchio/LG-Snippetbox/pkg/models"
)

// mockUser is mock user to be used in test, associated to ID=1
// in the backend
var mockUser = &models.User{
	ID:      1,
	Name:    "MockBoy",
	Email:   "mockboy@invalid.example",
	Created: time.Now(),
	Active:  true,
}

type UserModel struct{}

func (u *UserModel) Get(id int) (*models.User, error) {
	switch id {
	case 1:
		return mockUser, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (u *UserModel) Insert(name, email, password string) error {
	switch email {
	case "mockboy@invalid.example":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (u *UserModel) Authenticate(email, password string) (int, error) {
	switch email {
	case "mockboy@invalid.example":
		return 1, nil
	default:
		return 0, models.ErrInvalidCredentials
	}
}
