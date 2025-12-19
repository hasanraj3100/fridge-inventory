package utils

import "github.com/golang-jwt/jwt/v5"

type PayLoad struct {
	Sub      int    `json:"sub"` // Subject (user ID)
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	Secret []byte
}

func NewJWTManager(secret string) *JWTManager {
	return &JWTManager{
		Secret: []byte(secret),
	}
}

func (jm *JWTManager) GenerateToken(payload *PayLoad) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signedToken, err := token.SignedString(jm.Secret)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (jm *JWTManager) VerifyToken(tokenStr string) (*PayLoad, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &PayLoad{}, func(token *jwt.Token) (any, error) {
		return jm.Secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*PayLoad)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
