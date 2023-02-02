package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/nn", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/nn", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postData := url.Values{}
	postData.Add("a", "a")
	postData.Add("b", "a")
	postData.Add("c", "a")

	r, _ = http.NewRequest("POST", "/nn", nil)

	r.PostForm = postData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows does not have required fields when it does")
	}
}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/nn", nil)
	form := New(r.PostForm)

	has := form.Has("first_name")
	if has {
		t.Error("forms shows has field when it does not")
	}

	postData := url.Values{}
	postData.Add("first_name", "aa")
	form = New(postData)

	has = form.Has("first_name")
	if !has {
		t.Error("shows form does not have field when it does")
	}
}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/nn", nil)
	form := New(r.PostForm)

	form.MinLength("x", 3, r)
	if form.Valid() {
		t.Error("form shows min length for non-exist field")
	}

	isError := form.Errors.Get("x")
	if isError == "" {
		t.Error("should have an err, but did not get one")
	}

	postedValues := url.Values{}
	postedValues.Add("some_field", "nainainai")
	form = New(postedValues)

	form.MinLength("some_field", 10, r)
	if form.Valid() {
		t.Error("shows min length of 100 met is shorter")
	}

	postedValues = url.Values{}
	postedValues.Add("another_field", "abcd123")
	form = New(postedValues)

	form.MinLength("another_field", 1, r)
	if !form.Valid() {
		t.Error("show min length of 1 is not met when it is")
	}

	isError = form.Errors.Get("x2")
	if isError != "" {
		t.Error("should not have an err, but got one")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedValues := url.Values{}
	form := New(postedValues)

	form.IsEmail("x")
	if form.Valid() {
		t.Error("form shows valid email for non existent field")
	}

	postedValues = url.Values{}
	postedValues.Add("email", "abc@gmail.com")
	form = New(postedValues)

	form.IsEmail("email")
	if !form.Valid() {
		t.Error("got an invalid email when should not have")
	}

	postedValues = url.Values{}
	postedValues.Add("email", "x")
	form = New(postedValues)

	form.IsEmail("email")
	if form.Valid() {
		t.Error("got valid for invalid email address")
	}
}
