package mono_test

import (
	"reflect"
	"strings"
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

func expectErrorStartsWith(t *testing.T, actual error, startsWith string) {
	if actual == nil {
		t.Fatal("Expected error to be present, actual nil")
		return
	}

	if strings.Index(actual.Error(), startsWith) != 0 {
		t.Fatal("Actual error is '" + actual.Error() + "', expected to start with '" + startsWith + "'")
	}
}

func expectNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatal("Got error: " + err.Error())
	}
}

func expectEquals(t *testing.T, actual, expected interface{}) {
	if actual != expected {
		t.Fatalf("Actual != expected. Actual: '%v', expected: '%v'", actual, expected)
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
