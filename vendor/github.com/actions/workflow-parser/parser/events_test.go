package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsAllowedEventType(t *testing.T) {
	// This is not exhaustive. Rather, it's a demonstration of some of the interesting bits of the func.
	allowed := []string{
		"push",
		"PUSH",
		"pull_request",
	}

	for _, s := range allowed {
		assert.True(t, isAllowedEventType(s), "should allow %q", s)
	}

	// This is also not exhaustive. We want to have this done by universe, after all.
	notAllowed := []string{
		"installation",
		"randommashingofkeyboard",
	}

	for _, s := range notAllowed {
		assert.False(t, isAllowedEventType(s), "should not allow %q", s)
	}
}
