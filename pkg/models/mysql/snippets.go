/*package mysql containes the model for the connecition with a MySQL backend*/
package mysql

import (
	"database/sql"
	"errors"

	"github.com/Zicchio/LG-Snippetbox/pkg/models"
)

// type SnippetModel wraps a DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

// Insert a snippet in the DB. If succesfull, will return the id of the
// latest inserted snippet
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	// prepared statemenet (used to prevent SQL injections)
	// NOTE: the syntax of prepared statement depends on the Database driver
	statement := "INSERT INTO snippets (title, content, created, expires) VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))"
	res, err := m.DB.Exec(statement, title, content, expires)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	// NOTE: the following code can be shortened as err := m.DB.QueryRow("SELECT ...", id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	// fetch by id and only non expired snippets
	statement := "SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?"
	// we assume that up to 1 row can be found
	row := m.DB.QueryRow(statement, id)

	s := &models.Snippet{}

	// use scan function to parse the row into a struct
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// distinguish a not found error from from an internal error
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord // row not found error
		}
		return nil, err // generic error
	}
	return s, nil
}

// Get 10 most recently created snippets
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	statement := "SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10"
	rows, err := m.DB.Query(statement)
	if err != nil {
		return nil, err
	}
	// IMPORTANT: not closing the connection will eventually result in all the conncetions in the pool being used up
	defer rows.Close()
	snippets := make([]*models.Snippet, 0)
	for rows.Next() {
		s := &models.Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	// this will return true if ANY error was encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
