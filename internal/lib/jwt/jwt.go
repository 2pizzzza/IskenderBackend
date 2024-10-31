package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type User interface {
	GetID() int64
	GetName() string
	GetEmail() string
}

func NewToken(user User, jwtSecret string, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodES256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.GetID()
	claims["email"] = user.GetEmail()
	claims["username"] = user.GetName()
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

//func ValidateToken(ctx context.Context, jwtSecret string) (jwt.MapClaims, error) {
//	m, ok := metadata.FromIncomingContext(ctx)
//	if !ok {
//		return nil, fmt.Errorf("missing context")
//	}
//
//	authHead := m.Get("authorization")
//	if len(authHead) == 0 {
//		return nil, fmt.Errorf("missing auth header")
//	}
//
//	tokenStr := strings.TrimSpace(strings.TrimPrefix(authHead[0], "Bearer"))
//
//	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
//		return []byte(jwtSecret), nil
//	})
//
//	if err != nil {
//		return nil, fmt.Errorf("invalid token: %v", err)
//	}
//
//	if !token.Valid {
//		return nil, fmt.Errorf("invalid token")
//	}
//
//	claims, ok := token.Claims.(jwt.MapClaims)
//	if !ok {
//		return nil, fmt.Errorf("invalid claims")
//	}
//	return claims, nil
//}
