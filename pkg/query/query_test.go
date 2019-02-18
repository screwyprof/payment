package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAccountShortInfo_QueryID(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "GetAccountShortInfo", GetAccountShortInfo{}.QueryID())
}

func TestGetAllAccounts_QueryID(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "GetAllAccounts", GetAllAccounts{}.QueryID())
}

func TestUnknown_QueryID(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "Unknown", Unknown{}.QueryID())
}
