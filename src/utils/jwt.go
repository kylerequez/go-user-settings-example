package utils

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"

	"github.com/kylerequez/go-user-settings-example/src/models"
)

func GenerateJwtToken(user models.User) (string, error) {
	key, err := GetEnv("JWT_SECRET_KEY")
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"ID":        user.ID.String(),
		"Authority": user.Authority,
		"Theme":     user.Settings.Theme,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(*key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ExtractClaims(tokenString string) (jwt.MapClaims, error) {
	key, err := GetEnv("JWT_SECRET_KEY")
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("there was an error in parsing the jwt")
		}
		return []byte(*key), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("there was error in parsing the claims")
	}

	return claims, nil
}
