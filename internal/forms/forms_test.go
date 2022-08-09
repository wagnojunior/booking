package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when it should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData) //Add the post-form values of the request r to form

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when it should have been invalid")
	}

	postedData.Add("a", "value of a")
	postedData.Add("b", "value of b")
	postedData.Add("c", "value of c")

	form = New(postedData)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("form shows missing required field when it should not have")
	}
}

func TestForm_Has(t *testing.T) {
	// Create a test http request
	postedData := url.Values{}
	form := New(postedData)

	if form.Has("a") {
		t.Error("form shows it has a when it does not")
	}

	postedData.Add("a", "value of a")

	form = New(postedData)
	if !form.Has("a") {
		t.Error("form shows it does not have a when it does")
	}

}

func TestForm_MinLength(t *testing.T) {
	postedData := url.Values{}
	postedData.Add("a", "value of a")

	form := New(postedData) //Add the post-form values of the request r to form
	if form.MinLength("a", 50) {
		t.Error("form shows min length is satisfied when it is not")
	}

	isError := form.Errors.Get("a")
	if isError == "" {
		t.Error("should have an error but did not get one")
	}

	postedData.Add("b", "value of b")
	form = New(postedData) //Add the post-form values of the request r to form

	if !form.MinLength("b", 5) {
		t.Error("form shows min length is not satisfied when it is")
	}

	isError = form.Errors.Get("b")
	if isError != "" {
		t.Error("should not have an error but got one")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	postedData.Add("a", "a")

	form := New(postedData) //Add the post-form values of the request r to form
	form.IsEmail("a")

	if form.Valid() {
		t.Error("form shows valid email adress when it should be invalid")
	}

	postedData.Add("b", "b@gmail.com")

	form = New(postedData) //Add the post-form values of the request r to form
	form.IsEmail("b")

	if !form.Valid() {
		t.Error("form shows invalid email adress when it should be valid")
	}
}
