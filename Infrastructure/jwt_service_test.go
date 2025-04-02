package infrastructure

import (
    "os"
    "testing"
    "time"

    "github.com/golang-jwt/jwt"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
)

// JWTServiceTestSuite groups all JWT service-related tests
type JWTServiceTestSuite struct {
    suite.Suite
    mockSecret string
}

// SetupSuite runs once before all tests
func (suite *JWTServiceTestSuite) SetupSuite() {
    suite.mockSecret = "mock_secret"
    os.Setenv("JWT_SECRET", suite.mockSecret)
    jwtSecret = []byte(os.Getenv("JWT_SECRET"))
}

// TearDownSuite runs once after all tests
func (suite *JWTServiceTestSuite) TearDownSuite() {
    os.Unsetenv("JWT_SECRET")
}

// TestGenerateToken tests token generation
func (suite *JWTServiceTestSuite) TestGenerateToken() {
    userID := "12345"
    role := "user"

    token, err := GenerateToken(userID, role)
    assert.NoError(suite.T(), err)
    assert.NotEmpty(suite.T(), token)

    // Validate the generated token
    claims, err := ValidateToken(token)
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), userID, claims.UserID)
    assert.Equal(suite.T(), role, claims.Role)
}

// TestValidateToken_ValidToken tests validation of a valid token
func (suite *JWTServiceTestSuite) TestValidateToken_ValidToken() {
    userID := "12345"
    role := "admin"

    // Generate a valid token
    token, err := GenerateToken(userID, role)
    assert.NoError(suite.T(), err)

    // Validate the token
    claims, err := ValidateToken(token)
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), userID, claims.UserID)
    assert.Equal(suite.T(), role, claims.Role)
}

// TestValidateToken_InvalidToken tests validation of an invalid token
func (suite *JWTServiceTestSuite) TestValidateToken_InvalidToken() {
    invalidToken := "invalid.token.string"

    claims, err := ValidateToken(invalidToken)
    assert.Error(suite.T(), err)
    assert.Nil(suite.T(), claims)
}

// TestValidateToken_ExpiredToken tests validation of an expired token
func (suite *JWTServiceTestSuite) TestValidateToken_ExpiredToken() {
    // Create an expired token
    claims := Claims{
        UserID: "12345",
        Role:   "user",
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(-1 * time.Hour).Unix(), // Expired 1 hour ago
            IssuedAt:  time.Now().Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtSecret)
    assert.NoError(suite.T(), err)

    // Validate the expired token
    parsedClaims, err := ValidateToken(tokenString)
    assert.Error(suite.T(), err)
    assert.Nil(suite.T(), parsedClaims)
}

// TestValidateToken_TamperedToken tests validation of a tampered token
func (suite *JWTServiceTestSuite) TestValidateToken_TamperedToken() {
    userID := "12345"
    role := "user"

    // Generate a valid token
    token, err := GenerateToken(userID, role)
    assert.NoError(suite.T(), err)

    // Tamper with the token
    tamperedToken := token + "tampered"

    claims, err := ValidateToken(tamperedToken)
    assert.Error(suite.T(), err)
    assert.Nil(suite.T(), claims)
}

// Run the test suite
func TestJWTServiceTestSuite(t *testing.T) {
    suite.Run(t, new(JWTServiceTestSuite))
}
