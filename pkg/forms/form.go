package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Form is a general purpose struct for form data and errors associated to
// that data
type Form struct {
	url.Values // url.Values is embedded so that the template engine and our code can use Get
	Errors     errors
}

func New(data url.Values) *Form {
	return &Form{data, errors(map[string][]string{})}
}

// Required checks that the imput fields are in the form data
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be empty")
		}
	}
}

// MaxLenght checks that a text input field has at most d characters
func (f *Form) MaxLenght(field string, d int) {
	value := f.Get(field)
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (maximum length is %d letters)", d))
	}
}

func (f *Form) MinLength(field string, d int) {
	value := f.Get(field)
	if utf8.RuneCountInString(value) < d {
		f.Errors.Add(field, fmt.Sprintf("This field is too short (minimum length is %d letters)", d))
	}
}

// AdmittedValues checks that a field has a value in an enum of valid options
// opts
func (f *Form) AdmittedValues(field string, opts ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	for _, opt := range opts {
		if value == opt {
			return
		}
	}
	f.Errors.Add(field, "This field is invalid")
}

func (f *Form) MatchesPattern(field string, pat *regexp.Regexp) {
	value := f.Get(field)
	if value == "" {
		// empty field match all patterns (requirement check is perfomed elsewhere)
		return
	}
	if !pat.Match([]byte(value)) {
		f.Errors.Add(field, fmt.Sprintf("This field if invalid"))
	}
}

// Valid return true iff the form data is valid accordin to checks
// currently perfomed
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
