package middlewares

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"os"
)

// Authentication
// 업비트 JWT 인증토큰 발급방법
// https://docs.upbit.com/docs/create-authorization-request
func Authentication() string {
	accessKey := os.Getenv("AccessKey")
	secretKey := os.Getenv("SecretKey")

	payload := jwt.MapClaims{
		"access_key": accessKey,
		"nonce":      uuid.New().String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		fmt.Println("Error generating JWT token:", err)
		panic(err)
	}

	// Construct the authorization token
	authorizationToken := "Bearer " + tokenString
	fmt.Println("Authorization Token:", authorizationToken)
	return authorizationToken
}
