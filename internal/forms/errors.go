package forms

// errors is a custom type that maps a string to a slice of strings
// The map's key is a form field in <make-reservation.page.tmpl>
// The map's value is a set of messages associated with a key
type errors map[string][]string

// Add adds a message to a corresponding field in the variable e of type errors
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get gets a set of messages associated with a key and returns the first message
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}

	return es[0]
}
