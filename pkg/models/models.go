/*package models contains top level data types*/
package models

import (
	"errors"
	"time"
)

// ErrNoRecord is used as error inplace of sql.ErrNoRows in order to encapsulate the model completely
var ErrNoRecord = errors.New("models: no matching record found")

// SnippetDTO
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}
