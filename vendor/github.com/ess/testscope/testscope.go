package testscope

import (
	"os"
	"testing"
)

func SkipUnlessUnit(t *testing.T) {
	if !Unit() {
		t.Skip("Not running unit tests at this time")
	}
}

func SkipUnlessIntegration(t *testing.T) {
	if !Integration() {
		t.Skip("Not running integration tests at this time")
	}
}

func Unit() bool {
	if len(os.Getenv("UNIT")) > 0 {
		return true
	}

	return false
}

func Integration() bool {
	if len(os.Getenv("INTEGRATION")) > 0 {
		return true
	}

	return false
}
