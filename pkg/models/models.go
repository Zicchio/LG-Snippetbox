/*package models contains top level data types*/
package models

import (
	"errors"
	"time"
)

var (
	// ErrNoRecord is used as error inplace of sql.ErrNoRows in order to encapsulate the model completely
	ErrNoRecord = errors.New("models: no matching record found")
	// ErrInvalidCredentials is used if a user tries to login with invalid username or password
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	// ErrDuplicateEmail is used if a user signs up with an already existing email
	// Security NOTE: in some services, dyplaying to users that an email already exists might be considered a security issue
	ErrDuplicateEmail = errors.New("models: duplicate email")
)

// SnippetDTO
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// UserDTO
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword string
	Created        time.Time
	Active         bool
}
