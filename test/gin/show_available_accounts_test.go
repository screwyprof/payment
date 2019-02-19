package gin

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/screwyprof/payment/internal/pkg/delivery/gin/response"
)

func (s *GopherPayTestSuite) ShowAvailableAccounts(t *testing.T) {
	// arrange
	expected := []response.AvailableAccount{
		{
			Number: "ACC555",
		},
		{
			Number: "ACC777",
		},
	}

	// act
	obtained, err := s.client.ShowAvailableAccounts()
	require.NoError(t, err, "cannot show available accounts")

	//  slice elements' order is not guaranteed
	sort.Slice(obtained, func(i int, j int) bool { return i > j })

	// assert
	assert.Equal(t, expected, obtained)
}
