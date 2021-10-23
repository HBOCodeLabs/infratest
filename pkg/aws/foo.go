package aws

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertEqual(t *testing.T, expected string, actual string) {
	assert.Equal(t, expected, actual, fmt.Sprintf("expected '%s' to equal '%s", expected, actual))
}
