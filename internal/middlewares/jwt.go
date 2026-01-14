package middlewares

import (
	"crypto/sha512"
	"encoding/hex"
	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"upbit-api/config"
)

// CreateTokenWithNoParams 파라미터 없는 API용 JWT 토큰 생성
// https://docs.upbit.com/docs/create-authorization-request
func CreateTokenWithNoParams() string {
	payload := jwt.MapClaims{
		"access_key": config.AccessKey,
		"nonce":      uuid.New().String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		log.Println("JWT 토큰 생성 실패:", err)
		panic(err)
	}

	return "Bearer " + tokenString
}

// CreateTokenWithParams 파라미터 있는 API용 JWT 토큰 생성 (SHA512 해시 포함)
func CreateTokenWithParams(query string) string {
	hash := sha512.New()
	hash.Write([]byte(query))
	queryHash := hex.EncodeToString(hash.Sum(nil))

	payload := jwt.MapClaims{
		"access_key":     config.AccessKey,
		"nonce":          uuid.New().String(),
		"query_hash":     queryHash,
		"query_hash_alg": "SHA512",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		log.Println("JWT 토큰 생성 실패:", err)
		panic(err)
	}

	return "Bearer " + tokenString
}
