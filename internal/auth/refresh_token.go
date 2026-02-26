package auth

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type RefreshClaims struct {
	UserID int64 `json:"user_id"`
}

func ParseRefreshToken(tokenStr string) (*RefreshClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid user_id")
	}

	return &RefreshClaims{
		UserID: int64(userIDFloat),
	}, nil
}
