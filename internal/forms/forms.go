package forms

import (
	"net/http"
	"net/url"
)

// Form creates custom form struct, embeds an url.Values object
type Form struct {
	url.Values
	Errors errors
}

// Valid return true if there are no errs, otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// New initializes a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Has check if form field is in post and not empty
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		f.Errors.Add(field, "This field can't be blank")
		return false
	}
	return true
}