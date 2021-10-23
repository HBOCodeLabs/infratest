package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssertEqual(t *testing.T) {
	fakeTest := testing.T{}
	AssertEqual(t, "a", "b")

	assert.True(t, fakeTest.Failed())
}
