package auth

import (
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHasRole(t *testing.T) {

	// define echo context
	e := echo.New()
	c := e.AcquireContext()

	testedRole := "mocked_role"

	// Add role to claim
	claims := make(jwt.MapClaims)
	claims["http://supervisorapi/role"] = testedRole

	// define jwt token
	tk := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	// add jwt to context
	c.Set("user", tk)

	// call Has Role
	ok := HasRole(c, testedRole)

	// Assert result
	assert.True(t, ok)
}

func TestIsManager(t *testing.T) {

	// define echo context
	e := echo.New()
	c := e.AcquireContext()

	testedRole := "manager"

	// Add role to claim
	claims := make(jwt.MapClaims)
	claims["http://supervisorapi/role"] = testedRole

	// define jwt token
	tk := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	// add jwt to context
	c.Set("user", tk)

	// call Has Role
	ok := IsManager(c)

	// Assert result
	assert.True(t, ok)
}

func TestGetUserNickName(t *testing.T) {

	// define echo context
	e := echo.New()
	c := e.AcquireContext()

	// Add role to claim

	claims := make(jwt.MapClaims)
	claims["http://supervisorapi/nickname"] = "mocked_name"

	// define jwt token
	tk := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	// add jwt to context
	c.Set("user", tk)

	// call Has Role
	fetchedName := GetUserNickname(c)

	// Assert result
	assert.Equal(t, fetchedName, "mocked_name")
}

func TestGetUserId(t *testing.T) {

	// define echo context
	e := echo.New()
	c := e.AcquireContext()

	// Add role to claim

	claims := make(jwt.MapClaims)
	claims["sub"] = "mocked_id"

	// define jwt token
	tk := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	// add jwt to context
	c.Set("user", tk)

	// call Has Role
	fetchedId := GetUserId(c)

	// Assert result
	assert.Equal(t, fetchedId, "mocked_id")
}
