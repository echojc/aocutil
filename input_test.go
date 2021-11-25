package aocutil_test

import (
	"testing"

	"github.com/echojc/aocutil"
)

func TestInput(t *testing.T) {
	i, err := aocutil.NewInputFromFile("session_id")
	if err != nil {
		t.Errorf("Could not read session ID from file: %w", err)
		t.FailNow()
	}

	i.Ints(2018, 1)
}
