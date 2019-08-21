package common

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// ValidatorIntf interface for Validator
type ValidatorIntf interface {
	MustBeGreaterThan(high, value int) bool
	MustBeBefore(high, value time.Time) bool
	MustBeNotEmpty(fieldValue string) bool
	IsValid() bool
	Error() string
}

// Validator used for validation
type Validator struct {
	err []error
}

// NewValidator create a Validator struct
func NewValidator() *Validator {
	validator := Validator{}
	return &validator
}

// MustBeGreaterThan int comparison
func (v *Validator) MustBeGreaterThan(fieldName string, fieldvalue, high int) bool {
	if fieldvalue <= high {
		v.err = append(v.err, fmt.Errorf(fieldName+" Must be Greater than %d", high))
		return false
	}
	return true
}

// IsTimeBefore time comparison
func (v *Validator) IsTimeBefore(fieldName string, fieldvalue, max time.Time) bool {
	if fieldvalue.After(max) {
		v.err = append(v.err, fmt.Errorf(fieldName+"Must be Before than %v", max))
		return false
	}
	return true
}

// IsStrNotEmpty string not empty
func (v *Validator) IsStrNotEmpty(fieldName string, fieldValue string) bool {
	if fieldValue == "" {
		v.err = append(v.err, fmt.Errorf(fieldName+" Must not be Empty"))
		return false
	}
	return true
}

// IsInt64Negative returns true if value < 0
func (v *Validator) IsInt64Negative(fieldName string, fieldvalue int64) bool {
	if fieldvalue > 0 {
		v.err = append(v.err, fmt.Errorf(fieldName+" Must be Negative"))
		return false
	}
	return true
}

// IsInt64NonNegative returns true if value >= 0
func (v *Validator) IsInt64NonNegative(fieldName string, fieldvalue int64) bool {
	if fieldvalue <= 0 {
		v.err = append(v.err, fmt.Errorf(fieldName+" Must not be Negative"))
		return false
	}
	return true
}

// IsPositive returns true if value > 0
func (v *Validator) IsPositive(fieldName string, fieldvalue int64) bool {
	if fieldvalue < 0 {
		v.err = append(v.err, fmt.Errorf(fieldName+" Must be Positive"))
		return false
	}
	return true
}

// IsNonPositive returns true if value <= 0
func (v *Validator) IsNonPositive(fieldName string, fieldvalue int64) bool {
	if fieldvalue >= 0 {
		v.err = append(v.err, fmt.Errorf(fieldName+" Must not be Positive"))
		return false
	}
	return true
}

// IsEmail validate email address
func (v *Validator) IsEmail(fieldName string, fieldValue string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !re.MatchString(fieldValue) {
		v.err = append(v.err, fmt.Errorf(fieldName+" Not valid"))
		return false
	}
	return true
}

// IsDateFormat validate date format
func (v *Validator) IsDateFormat(fieldName string, fieldValue string) bool {
	re := regexp.MustCompile("(0?[1-9]|[12][0-9]|3[01])/(0?[1-9]|1[012])/((19|20)\\d\\d)")
	if !re.MatchString(fieldValue) {
		v.err = append(v.err, fmt.Errorf("Date format not valid for "+fieldName))
		return false
	}
	return true
}

// IsPhoneNumber validate phone number
func (v *Validator) IsPhoneNumber(fieldName string, fieldValue string) bool {
	re := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
	if !re.MatchString(fieldValue) {
		v.err = append(v.err, fmt.Errorf(fieldName+" Not valid"))
		return false
	}
	return true
}

// IsUUID4 validate UUID4
func (v *Validator) IsUUID4(fieldName string, fieldValue string) bool {
	re := regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$")
	if !re.MatchString(fieldValue) {
		v.err = append(v.err, fmt.Errorf("UUID not valid for "+fieldName))
		return false
	}
	return true
}

// IsAlpha validate Alpha characters
func (v *Validator) IsAlpha(fieldName string, fieldValue string) bool {
	re := regexp.MustCompile("^[a-zA-Z]+$")
	if !re.MatchString(fieldValue) {
		v.err = append(v.err, fmt.Errorf(fieldName+"not valid"))
		return false
	}
	return true
}

// IsAlphaNumeric validate alpha numeric characters
func (v *Validator) IsAlphaNumeric(fieldName string, fieldValue string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9]+$")
	if !re.MatchString(fieldValue) {
		v.err = append(v.err, fmt.Errorf("alpha numeric characters not valid "+fieldName))
		return false
	}
	return true
}

// IsDigits validate Digits
func (v *Validator) IsDigits(fieldName string, fieldValue string) bool {
	re := regexp.MustCompile("^[+-]?([0-9]*\\.?[0-9]+|[0-9]+\\.?[0-9]*)([eE][+-]?[0-9]+)?$")
	if !re.MatchString(fieldValue) {
		v.err = append(v.err, fmt.Errorf("Digits not valid for "+fieldName))
		return false
	}
	return true
}

// IsCreditCard validate credit card
func (v *Validator) IsCreditCard(fieldName string, fieldValue string) bool {
	re := regexp.MustCompile("^(?:4[0-9]{12}(?:[0-9]{3})?|[25][1-7][0-9]{14}|6(?:011|5[0-9][0-9])[0-9]{12}|3[47][0-9]{13}|3(?:0[0-5]|[68][0-9])[0-9]{11}|(?:2131|1800|35\\d{3})\\d{11})$")
	if !re.MatchString(fieldValue) {
		v.err = append(v.err, fmt.Errorf("Details not valid for "+fieldName))
		return false
	}
	return true
}

// IsStrLenGtMin string length greater than min
func (v *Validator) IsStrLenGtMin(fieldName string, fieldValue string, min int) bool {

	if len(fieldValue) < min {
		v.err = append(v.err, fmt.Errorf("Number of characters of "+fieldName+" must be greater than %d", min))
		return false
	}

	return true
}

// IsStrLenLtMax string length less than max
func (v *Validator) IsStrLenLtMax(fieldName string, fieldValue string, max int) bool {

	if len(fieldValue) > max {
		v.err = append(v.err, fmt.Errorf("Number of characters of "+fieldName+" must be less than %d", max))
		return false
	}

	return true
}

// IsStrLenBetMinMax string length between min max
func (v *Validator) IsStrLenBetMinMax(fieldName string, fieldValue string, min int, max int) bool {

	if (len(fieldValue) < min) || (len(fieldValue) > max) {
		v.err = append(v.err, fmt.Errorf("Number of characters of "+fieldName+" must be between %d and %d", min, max))
		return false
	}

	return true
}

// IsValid bool check
func (v *Validator) IsValid() bool {
	return v.err != nil
}

// Error error
func (v *Validator) Error() string {
	var x []string
	for _, err := range v.err {
		x = append(x, err.Error())
	}
	return strings.Join(x, ", ")
}
