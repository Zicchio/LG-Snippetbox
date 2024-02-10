package forms

// custom "error" type for validation errros - not in the sense of errors package
// Personal NOTE: should have been named validationErrors or something
type errors map[string][]string

// Add adds a value to the error map
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get returns the first error for a given field, or empty string if no errors
// for that field were obtained
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0] // NOTE: why not concant?
}
