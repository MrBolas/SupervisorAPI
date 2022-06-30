package auth

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type jwtCustomClaims struct {
	Username string `json:"preferred_username"`
	jwt.StandardClaims
}

type KeyResponse struct {
	PublicKeys []Aoth0Key `json:"keys"`
}

type Aoth0Key struct {
	Alg string
	Kty string
	Use string
	N   string
	E   string
	Kid string
	X5t string
	X5c []string
}

func JWTConfig(publicKeyUrl string) (middleware.JWTConfig, error) {

	publicKey, err := getPublicKey(publicKeyUrl)
	if err != nil {
		return middleware.JWTConfig{}, err
	}

	pemKey := generatePEMKey(publicKey)

	k, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pemKey))
	if err != nil {
		return middleware.JWTConfig{}, err
	}
	jwtConfig := middleware.JWTConfig{
		SigningMethod: "RS256",
		SigningKey:    k,
	}

	return jwtConfig, nil
}

func generatePEMKey(publicKey string) string {

	lines := int(math.Ceil(float64(len(publicKey)) / 64))
	finalKey := "-----BEGIN PUBLIC KEY-----\n"
	for i := 0; i < lines; i++ {
		if i != lines-1 {
			finalKey += publicKey[i*64:i*64+64] + "\n"
		} else { // last line
			finalKey += publicKey[i*64:] + "\n"
		}
	}
	finalKey += "-----END PUBLIC KEY-----"

	return finalKey
}

func getPublicKey(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var keyResponse KeyResponse
	err = json.Unmarshal(bodyBytes, &keyResponse)
	if err != nil {
		return "", err
	}

	return keyResponse.PublicKeys[0].X5c[0], nil
}

func HasRole(c echo.Context, permittedRole string) bool {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	role := claims["http://supervisorapi/role"].(string)

	return role == permittedRole
}

func IsManager(c echo.Context) bool {
	return HasRole(c, "manager")
}

func GetUserNickname(c echo.Context) string {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	nickname := claims["http://supervisorapi/nickname"].(string)
	return nickname
}

func GetUserId(c echo.Context) string {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["sub"].(string)
	return userId
}
