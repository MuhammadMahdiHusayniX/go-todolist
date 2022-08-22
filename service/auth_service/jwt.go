package auth_service

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTService interface {
	GenerateToken(email string, isUser bool) string
	ValidateToken(token string) (*jwt.Token, error)
}

type authCustomClaims struct {
	Name string `json:"name"`
	User bool   `json:"user"`
	jwt.RegisteredClaims
}

type jwtServices struct {
	secretKey string
	issure    string
}

// auth-jwt
func JWTAuthService() JWTService {
	return &jwtServices{
		secretKey: getSecretKey(),
		issure:    "Mahdi",
	}
}

func getSecretKey() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret"
	}
	return secret
}

func (service *jwtServices) GenerateToken(email string, isUser bool) string {
	claims := &authCustomClaims{
		email,
		isUser,
		jwt.RegisteredClaims{
			Issuer:    service.issure,
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour * 48)},
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
			ID:        email,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//encoded string
	t, err := token.SignedString([]byte(service.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

func (service *jwtServices) ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, fmt.Errorf("invalid token")
		}
		return []byte(service.secretKey), nil
	})

}
