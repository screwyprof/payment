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

func (s *GopherPayTestSuite) OpenAccount_OpenACC777(t *testing.T) {
	// arrange
	q := request.OpenAccount{
		Number:   "ACC777",
		Amount:   77700,
		Currency: "USD",
	}

	expected := response.ShortAccountInfo{
		Number:  "ACC777",
		Balance: "$777.00",
	}

	// act
	obtained, err := s.client.OpenAccount(q)
	require.NoError(t, err, "cannot open account")

	// assert
	assert.Equal(t, expected, obtained)
}

func (s *GopherPayTestSuite) OpenAccount_OpenACC555(t *testing.T) {
	// arrange
	q := request.OpenAccount{
		Number:   "ACC555",
		Amount:   55500,
		Currency: "USD",
	}

	expected := response.ShortAccountInfo{
		Number:  "ACC555",
		Balance: "$555.00",
	}

	// act
	obtained, err := s.client.OpenAccount(q)
	require.NoError(t, err, "cannot open account")

	// assert
	assert.Equal(t, expected, obtained)
}
