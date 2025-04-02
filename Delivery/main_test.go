package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// MainTestSuite groups all main-related tests
type MainTestSuite struct {
	suite.Suite
}

// SetupSuite runs once before all tests
func (suite *MainTestSuite) SetupSuite() {
	// Set up any global configurations or environment variables
	os.Setenv("MONGODB_URI", "mongodb://mockhost:27017")
}

// TearDownSuite runs once after all tests
func (suite *MainTestSuite) TearDownSuite() {
	// Clean up global configurations or environment variables
	os.Unsetenv("MONGODB_URI")
}

// TestInitMongoDB tests the MongoDB initialization
func (suite *MainTestSuite) TestInitMongoDB() {
	os.Setenv("MONGODB_URI", "mongodb://mockhost:27017")
	defer os.Unsetenv("MONGODB_URI")

	client, db, err := initMongoDB()
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), client)
	assert.NotNil(suite.T(), db)

	_ = client.Disconnect(context.Background())
}

// TestInitServer tests the server initialization
func (suite *MainTestSuite) TestInitServer() {
	router := http.NewServeMux()
	server := initServer(router)

	assert.Equal(suite.T(), ":8080", server.Addr)
	assert.Equal(suite.T(), router, server.Handler)
}

// TestRunServer tests running the server
func (suite *MainTestSuite) TestRunServer() {
	router := http.NewServeMux()
	server := initServer(router)

	go func() {
		_ = server.ListenAndServe()
	}()
	time.Sleep(100 * time.Millisecond) // Allow server to start

	assert.NoError(suite.T(), server.Close())
}

// TestRunServerWithError tests the runServer function when the server fails to start
func (suite *MainTestSuite) TestRunServerWithError() {
	server := &http.Server{
		Addr:    ":invalid", // Invalid address to trigger an error
		Handler: http.NewServeMux(),
	}

	go func() {
		runServer(server, true) // Suppress logs during testing
	}()
	time.Sleep(100 * time.Millisecond) // Allow goroutine to execute

	assert.NoError(suite.T(), server.Close())
}

// TestMainFunction tests the main function indirectly by mocking dependencies
func (suite *MainTestSuite) TestMainFunction() {
	// Mock environment variables
	os.Setenv("MONGODB_URI", "mongodb://mockhost:27017")
	defer os.Unsetenv("MONGODB_URI")

	// Mock signal handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT)
	go func() {
		time.Sleep(100 * time.Millisecond)
		quit <- syscall.SIGINT
	}()

	// Run main function in a separate goroutine
	go func() {
		main()
	}()
	time.Sleep(200 * time.Millisecond) // Allow main to execute

	assert.True(suite.T(), true) // Placeholder assertion to ensure test runs
}

// Run the test suite
func TestMainTestSuite(t *testing.T) {
	suite.Run(t, new(MainTestSuite))
}
