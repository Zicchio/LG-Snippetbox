package mock

import (
	"time"

	"github.com/Zicchio/LG-Snippetbox/pkg/models"
)

// mockSnippet to be used in tests associated to id=1 in DB
var mockSnippet = &models.Snippet{
	ID:      1,
	Title:   "Snippet tile here",
	Content: "Snippet content goes here...",
	Created: time.Now(),
	Expires: time.Now(),
}

type SnippetModel struct{}

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	return 2, nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}
