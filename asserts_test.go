package mono_test

import (
	"reflect"
	"testing"
)

func expectNotNil(t *testing.T, v interface{}) {
	if v == nil {
		t.Fatal("expected value or pointer, got nil")
	}
}

func expectError(t *testing.T, actual error, expected string) {
	if actual == nil {
		t.Fatal("Expected error to be present, actual nil")
		return
	}

	if actual.Error() != expected {
		t.Fatal("Actual error is '" + actual.Error() + "', expected '" + expected + "'")
	}
}

func expectNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatal("Got error: " + err.Error())
	}
}

func expectDeepEquals(t *testing.T, actual, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		t.Fatal("Actual != expected")
	}
}

func expectTrue(t *testing.T, v bool) {
	if !v {
		t.Fatal("Got false, expected true")
	}
}
