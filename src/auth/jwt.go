package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey []byte

type JWTClaim struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	UserId   string `json:"user_id"`
	jwt.RegisteredClaims
}

func ValidateJWTToken(signedToken string) error {
	// Parse token
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return err
	}

	// Check if token is valid
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return ErrFailedToParseClaims
	}
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return ErrTokenExpired
	}

	return nil
}

func GetUserIdFromToken(signedToken string) (string, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return "", err
	}

	// Check if token is valid
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return "", ErrFailedToParseClaims
	}
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return "", ErrTokenExpired
	}

	return claims.UserId, nil
}

func GenerateJWTToken(username string, email string, userId string) (string, error) {
	expirationTime := time.Now().Add(12 * time.Hour)

	claims := &JWTClaim{
		Username: username,
		Email:    email,
		UserId:   userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: expirationTime,
			},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func Init() {
	jwtKey = []byte(os.Getenv("JWT_SECRET"))

	// I do not know if this is the best place to put this
	// but it is the only place I can think of
	CreateAdminUser()
}
