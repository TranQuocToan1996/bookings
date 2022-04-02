package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Hold form information
type Form struct {
	url.Values
	Errors errors
}

// New initialize a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Has if form field is in post request and not empty value
func (f *Form) Has(field string) bool {
	valueField := f.Get(field)
	return valueField != ""
}

// Check for required fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field can't be blank")
		}
	}
}

// Check the form from receiver valid or not
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// MinLength check for string minimum length
func (f *Form) MinLength(field string, length int) bool {
	value := f.Get(field)
	if len(value) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d character long", length))
		return false
	}
	return true
}

// Check for valid email address
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address!")
	}
}

// Check for valid phone number
func (f *Form) IsPhoneNumber(field string) {
	// This pattern use for VN, US area
	regexpPatternListPhoneNum := []string{
		`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`,
		`(84|0[3|5|7|8|9])+([0-9]{8})\b`,
		`^[0-9\-\+]{9,15}$`,
		`^(\+\d{1,3}[- ]?)?\d{10}$`,
		`(84|0[3|5|7|8|9])+([0-9]{8})\b`,
	}

	for _, pattern := range regexpPatternListPhoneNum {
		re := regexp.MustCompile(pattern)
		if re.MatchString(f.Get(field)) {
			return
		}
	}
	f.Errors.Add(field, "Invalid Phone Number!")
}
