package gin

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/screwyprof/payment/test/app"
)

// TestGopherPayTestSuite initialize test suit.
func TestGopherPayTestSuite(t *testing.T) {
	suite.Run(t, new(GopherPayTestSuite))
}

// GopherPayTestSuite tests GopherPay service.
type GopherPayTestSuite struct {
	app.AppTestSuite

	client *GopherPayTestClient
}

// SetupSuite runs once on suit initialization.
func (s *GopherPayTestSuite) SetupSuite() {
	// call parent
	s.AppTestSuite.SetupSuite()

	s.client = &GopherPayTestClient{
		client: &HttpTestClient{
			handler: s.Router,
		},
	}
}

func (s *GopherPayTestSuite) TestGopherPay() {
	s.T().Run("OpenAccount", func(t *testing.T) {
		// Test creating ACC777 with balance $777.00
		s.OpenAccount_OpenACC777(t)
		// Test creating ACC555 with balance $555.00
		s.OpenAccount_OpenACC555(t)
	})

	s.T().Run("TransferMoney", func(t *testing.T) {
		// transfer money from ACC777 to ACC555
		s.TransferMoney(t)
	})

	s.T().Run("ShowAccount", func(t *testing.T) {
		s.ShowAccount_ShowACC777(t)
	})
}
