package infrastructure

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// AuthMiddlewareTestSuite groups all middleware-related tests
type AuthMiddlewareTestSuite struct {
	suite.Suite
	router *gin.Engine
}

// SetupSuite runs once before all tests
func (suite *AuthMiddlewareTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

// SetupTest runs before each test
func (suite *AuthMiddlewareTestSuite) SetupTest() {
	suite.router = gin.New()
}

// TestAuthMiddleware_MissingAuthorizationHeader tests missing Authorization header
func (suite *AuthMiddlewareTestSuite) TestAuthMiddleware_MissingAuthorizationHeader() {
	suite.router.Use(AuthMiddleware(func(token string) (*Claims, error) {
		return nil, errors.New("mock validation not implemented")
	}))
	suite.router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, resp.Code)
	assert.JSONEq(suite.T(), `{"error": "authorization header is required"}`, resp.Body.String())
}

// TestAuthMiddleware_InvalidAuthorizationHeaderFormat tests invalid Authorization header format
func (suite *AuthMiddlewareTestSuite) TestAuthMiddleware_InvalidAuthorizationHeaderFormat() {
	suite.router.Use(AuthMiddleware(func(token string) (*Claims, error) {
		return nil, errors.New("mock validation not implemented")
	}))
	suite.router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "InvalidHeader")
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, resp.Code)
	assert.JSONEq(suite.T(), `{"error": "invalid authorization header format"}`, resp.Body.String())
}

// TestAuthMiddleware_InvalidToken tests invalid token
func (suite *AuthMiddlewareTestSuite) TestAuthMiddleware_InvalidToken() {
	suite.router.Use(AuthMiddleware(func(token string) (*Claims, error) {
		return nil, errors.New("mock validation not implemented")
	}))
	suite.router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, resp.Code)
	assert.JSONEq(suite.T(), `{"error": "invalid token"}`, resp.Body.String())
}

// TestAuthMiddleware_ValidToken tests valid token
func (suite *AuthMiddlewareTestSuite) TestAuthMiddleware_ValidToken() {
	validToken := "valid_token"
	mockValidateToken := func(token string) (*Claims, error) {
		return &Claims{UserID: "123", Role: "user"}, nil
	}

	suite.router.Use(AuthMiddleware(mockValidateToken))
	suite.router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+validToken)
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	assert.JSONEq(suite.T(), `{"message": "success"}`, resp.Body.String())
}

// TestAuthMiddleware_ExpiredToken tests expired token
func (suite *AuthMiddlewareTestSuite) TestAuthMiddleware_ExpiredToken() {
	expiredToken := "expired_token"
	suite.router.Use(AuthMiddleware(func(token string) (*Claims, error) {
		return nil, errors.New("token expired")
	}))
	suite.router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+expiredToken)
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, resp.Code)
	assert.JSONEq(suite.T(), `{"error": "invalid token"}`, resp.Body.String())
}

// TestAdminMiddleware_NonAdminUser tests non-admin user access
func (suite *AuthMiddlewareTestSuite) TestAdminMiddleware_NonAdminUser() {
	suite.router.Use(func(c *gin.Context) {
		c.Set("role", "user")
		c.Next()
	})
	suite.router.Use(AdminMiddleware())
	suite.router.GET("/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest(http.MethodGet, "/admin", nil)
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusForbidden, resp.Code)
	assert.JSONEq(suite.T(), `{"error": "admin access required"}`, resp.Body.String())
}

// TestAdminMiddleware_AdminUser tests admin user access
func (suite *AuthMiddlewareTestSuite) TestAdminMiddleware_AdminUser() {
	suite.router.Use(func(c *gin.Context) {
		c.Set("role", "admin")
		c.Next()
	})
	suite.router.Use(AdminMiddleware())
	suite.router.GET("/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest(http.MethodGet, "/admin", nil)
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	assert.JSONEq(suite.T(), `{"message": "success"}`, resp.Body.String())
}

// Run the test suite
func TestAuthMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(AuthMiddlewareTestSuite))
}
