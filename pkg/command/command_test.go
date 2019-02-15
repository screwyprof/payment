package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenAccount_CommandID(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "OpenAccount", OpenAccount{}.CommandID())
}

func TestUnknown_CommandID(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "Unknown", Unknown{}.CommandID())
}
