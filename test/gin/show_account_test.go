package gin

import (
	//	"github.com/stretchr/testify/assert"
	//	"github.com/stretchr/testify/require"

	"github.com/screwyprof/payment/internal/pkg/delivery/gin/response"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/stretchr/testify/require"
)

func (s *GopherPayTestSuite) ShowAccount_ShowACC777(t *testing.T) {
	// arrange
	expected := response.AccountInfo{
		Number:  "ACC777",
		Balance: "$677.00",
		Ledgers: []response.Ledger{
			{Action: "Deposit, $777.00"},
			{Action: "Transfer to ACC555, $100.00"},
		},
	}

	// act
	obtained, err := s.client.ShowAccount("ACC777")
	require.NoError(t, err, "cannot show account")

	// assert
	assert.Equal(t, expected, obtained)
}
