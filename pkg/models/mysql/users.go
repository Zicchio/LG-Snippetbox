package mysql

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/Zicchio/LG-Snippetbox/pkg/models"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

// Insert adds a new record to the database
func (m *UserModel) Insert(name, email, password string) error {
	hashedPswd, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (name, email, hashed_password, created) VALUES (?,?,?, UTC_TIMESTAMP())`
	_, err = m.DB.Exec(stmt, name, email, hashedPswd)
	if err != nil {
		// if mysql error, check if it is due to dubplicate mail and in case add some context
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

// Authenticate verificy if a user (and corresponding credentials exists)
// If it exists, user id is returned
func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPswd []byte
	stmt := `SELECT id, hashed_password FROM users WHERE email= ? AND active=TRUE`
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hashedPswd)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		}
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(hashedPswd, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		}
		return 0, err
	}
	return id, nil
}

// Get fetches details for a specific given user
func (m *UserModel) Get(id int) (*models.User, error) {
	u := &models.User{}
	stmt := `SELECT id, name, email, created, active FROM users WHERE id = ?`
	err := m.DB.QueryRow(stmt, id).Scan(&u.ID, &u.Name, &u.Email, &u.Created, &u.Active)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return u, nil
}
