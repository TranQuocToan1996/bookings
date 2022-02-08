package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	if !form.Valid() {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	r, _ = http.NewRequest("POST", "/whatever", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows does not have required fields when it does")
	}
}

func TestForm_Has(t *testing.T) {

	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)
	if form.Has("whateverfield") {
		t.Error("Form have value when it shouldn't")
	}

	postedData := url.Values{}
	postedData.Add("example", "example")
	form = New(postedData)
	if !form.Has("example") {
		t.Error("Field don't have value when it should")
	}
}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	// 10 just random number non-zero
	form.MinLength("x", 10)
	if form.Valid() {
		t.Error("Form shows Minlength for non-exist field")
	}

	// Test for errors.go
	isError := form.Errors.Get("x")
	if isError == "" {
		t.Error("should have an error, but did not get")
	}

	postedData := url.Values{}
	postedData.Add("example", "exampleVal")
	form = New(postedData)
	form.MinLength("example", 100)
	if form.Valid() {
		t.Error(`The value of string "exampleVal" is valid even is shorter than 100 characters`)
	}

	postedData = url.Values{}
	postedData.Add("example", "exampleVal")
	form = New(postedData)
	form.MinLength("example", 1)
	if !form.Valid() {
		t.Error(`The value of string "exampleVal" is invalid even it has more than 1 character`)
	}

	// Test for errors.go
	isError = form.Errors.Get("example")
	if isError != "" {
		t.Error("shouldn't have an error, but did get")
	}

}

func TestForm_IsEmail(t *testing.T) {
	/*
		r := httptest.NewRequest("POST", "/whatever", nil)
		form := New(r.PostForm)
	Old code */
	postedData := url.Values{}
	form := New(postedData)

	// Check for a field dont exist
	form.IsEmail("whateverfield")
	if form.Valid() {
		t.Error("non-exist field passing test")
	}

	// Check non-pass IsEmail for valid email address
	postedData = url.Values{}
	postedData.Add("email", "example@example.com")
	form = New(postedData)
	form.IsEmail("email")
	if !form.Valid() {
		t.Error("return t.Error for a valid email")
	}

	// Check pass IsEmail for invalid email address
	postedData = url.Values{}
	postedData.Add("email", "example")
	form = New(postedData)
	form.IsEmail("email")
	if form.Valid() {
		t.Error("Do not return t.Error for an invalid email")
	}

}

type dataForTestIsPhoneNumber struct {
	valid   []string
	invalid []string
}

func TestForm_IsPhoneNumber(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	var testdata = dataForTestIsPhoneNumber{
		invalid: []string{
			"!@#$%^&*(()",
			"abcdef",
		},
		valid: []string{
			"0999999999",
			"+84999999999",
			"(+84)999999999",
		},
	}
	// Check for a field dont exist
	form.IsPhoneNumber("whateverfield")
	if !form.Valid() {
		t.Error("non-exist field passing test")
	}

	// Check case for invalid phone
	for _, val := range testdata.invalid {
		postedData = url.Values{}
		postedData.Add("phone", val)
		form = New(postedData)
		form.IsPhoneNumber("phone")
		if form.Valid() {
			t.Errorf(`Invalid phone number "%s" is considered as valid`, val)
		}
	}

	// Check case for valid phone
	for _, val := range testdata.valid {
		postedData = url.Values{}
		postedData.Add("phone", val)
		form = New(postedData)
		form.IsPhoneNumber("phone")
		if !form.Valid() {
			t.Errorf(`Valid phone number "%s" is considered as invalid`, val)
		}
	}

}
