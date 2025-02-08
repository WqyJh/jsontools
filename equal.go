package jsontools

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/stretchr/testify/assert"
)

func JsonEqual(a, b []byte) (bool, error) {
	filter := NewJsonNullFilter(false)
	a, err := filter.Filter(a)
	if err != nil {
		return false, err
	}
	b, err = filter.Filter(b)
	if err != nil {
		return false, err
	}

	var o1, o2 any
	err = json.Unmarshal(a, &o1)
	if err != nil {
		return false, err
	}
	err = json.Unmarshal(b, &o2)
	if err != nil {
		return false, err
	}

	return reflect.DeepEqual(o1, o2), nil
}

type tHelper interface {
	Helper()
}

// AssertJSONEq asserts that two JSON strings are equivalent.
//
//	AssertJSONEq(t, `{"hello": "world", "foo": "bar"}`, `{"foo": "bar", "hello": "world"}`)
func AssertJSONEq(t assert.TestingT, expected string, actual string, msgAndArgs ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	filter := NewJsonNullFilter(true)
	expectedBytes, err := filter.Filter([]byte(expected))
	if err != nil {
		return assert.Fail(t, fmt.Sprintf("Expected value ('%s') is not valid json.\nJSON parsing error: '%s'", expected, err.Error()), msgAndArgs...)
	}
	actualBytes, err := filter.Filter([]byte(actual))
	if err != nil {
		return assert.Fail(t, fmt.Sprintf("Input ('%s') needs to be valid json.\nJSON parsing error: '%s'", actual, err.Error()), msgAndArgs...)
	}
	var expectedJSONAsInterface, actualJSONAsInterface interface{}

	if err := json.Unmarshal(expectedBytes, &expectedJSONAsInterface); err != nil {
		return assert.Fail(t, fmt.Sprintf("Expected value ('%s') is not valid json.\nJSON parsing error: '%s'", expected, err.Error()), msgAndArgs...)
	}

	if err := json.Unmarshal(actualBytes, &actualJSONAsInterface); err != nil {
		return assert.Fail(t, fmt.Sprintf("Input ('%s') needs to be valid json.\nJSON parsing error: '%s'", actual, err.Error()), msgAndArgs...)
	}

	return assert.Equal(t, expectedJSONAsInterface, actualJSONAsInterface, msgAndArgs...)
}

// RequireJSONEq asserts that two JSON strings are equivalent.
//
//	RequireJSONEq(t, `{"hello": "world", "foo": "bar"}`, `{"foo": "bar", "hello": "world"}`)
func RequireJSONEq(t assert.TestingT, expected string, actual string, msgAndArgs ...interface{}) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if AssertJSONEq(t, expected, actual, msgAndArgs...) {
		return
	}
	assert.FailNow(t, "JSON strings are not equal", msgAndArgs...)
}
