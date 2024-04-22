package token

import (
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

type SignedDetails struct {
	Email    string
	UserName string
	jwt.RegisteredClaims
}

var SECRET_KEY = os.Getenv("SECRET_KEY")

// Function to generate a signed JWT token from user details
func TokenGenerator(email string, username string) (signedtoken string, err error) {
	claims := &SignedDetails{
		Email:    email,
		UserName: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * time.Duration(24))),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}

	return token, err
}

// Function to validate a JWT token and return the user claims
func ValidateToken(signedtoken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(signedtoken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		msg = err.Error()
		return
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "The Token is invalid"
		return
	}

	timeValue := (*claims.ExpiresAt).Time

	if timeValue.Unix() < time.Now().Local().Unix() {
		msg = "token is expired"
		return
	}
	return claims, msg
}
