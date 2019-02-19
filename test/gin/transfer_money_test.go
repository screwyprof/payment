package gin

import (
	//	"github.com/stretchr/testify/assert"
	//	"github.com/stretchr/testify/require"

	"github.com/screwyprof/payment/internal/pkg/delivery/gin/response"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/screwyprof/payment/internal/pkg/delivery/gin/request"
)

func (s *GopherPayTestSuite) TransferMoney(t *testing.T) {
	// arrange
	q := request.TransferMoney{
		From:     "ACC777",
		To:       "ACC555",
		Amount:   10000,
		Currency: "USD",
	}

	expected := response.Message{
		Message: "$100.00 transfered from ACC777 to ACC555",
	}

	// act
	obtained, err := s.client.TransferMoney(q)
	require.NoError(t, err, "cannot transfer money")

	// assert
	assert.Equal(t, expected, obtained)
}
