package forms

type errors map[string][]string

// Add adds an error for a given name-attribute of form input
func (e errors) Add(field, msg string) {
	e[field] = append(e[field], msg)
}

// Get returns the first error message
func (e errors) Get(field string) string {
	errorString := e[field]
	if len(errorString) == 0 {
		return ""
	}
	return errorString[0]
}

