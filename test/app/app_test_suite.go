package app

import (
	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/suite"

	"github.com/screwyprof/payment/internal/pkg/app"
)

// AppTestSuite defines a test suite for the application.
type AppTestSuite struct {
	suite.Suite
	Router *gin.Engine
}

// SetupSuite runs once on suit initialization.
func (s *AppTestSuite) SetupSuite() {
	s.Router = app.SetupRouter()
}
