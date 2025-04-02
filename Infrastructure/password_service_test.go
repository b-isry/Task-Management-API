package infrastructure

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// PasswordServiceTestSuite groups all password service-related tests
type PasswordServiceTestSuite struct {
	suite.Suite
	password string
}

// SetupSuite runs once before all tests
func (suite *PasswordServiceTestSuite) SetupSuite() {
	suite.password = "securepassword123"
}

// TestHashPassword tests the HashPassword function
func (suite *PasswordServiceTestSuite) TestHashPassword() {
	// Hash the password
	hashedPassword, err := HashPassword(suite.password)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), hashedPassword)

	// Ensure the hashed password is not the same as the plain password
	assert.NotEqual(suite.T(), suite.password, hashedPassword)
}

// TestComparePasswords_ValidPassword tests ComparePasswords with a valid password
func (suite *PasswordServiceTestSuite) TestComparePasswords_ValidPassword() {
	// Hash the password
	hashedPassword, err := HashPassword(suite.password)
	assert.NoError(suite.T(), err)

	// Compare the hashed password with the correct plain password
	isValid := ComparePasswords(hashedPassword, suite.password)
	assert.True(suite.T(), isValid)
}

// TestComparePasswords_InvalidPassword tests ComparePasswords with an invalid password
func (suite *PasswordServiceTestSuite) TestComparePasswords_InvalidPassword() {
	invalidPassword := "wrongpassword"

	// Hash the password
	hashedPassword, err := HashPassword(suite.password)
	assert.NoError(suite.T(), err)

	// Compare the hashed password with an incorrect plain password
	isValid := ComparePasswords(hashedPassword, invalidPassword)
	assert.False(suite.T(), isValid)
}

// TestHashPassword_UniqueHashes tests that HashPassword generates unique hashes
func (suite *PasswordServiceTestSuite) TestHashPassword_UniqueHashes() {
	// Hash the same password multiple times
	hashedPassword1, err1 := HashPassword(suite.password)
	hashedPassword2, err2 := HashPassword(suite.password)

	assert.NoError(suite.T(), err1)
	assert.NoError(suite.T(), err2)

	// Ensure the hashes are different due to bcrypt's salting mechanism
	assert.NotEqual(suite.T(), hashedPassword1, hashedPassword2)
}

// Run the test suite
func TestPasswordServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PasswordServiceTestSuite))
}
