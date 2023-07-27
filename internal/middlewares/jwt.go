package middlewares

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"io"
	"net/http"
	"net/url"
	"upbit-api/config"
)

// 업비트 JWT 인증토큰 발급방법
// https://docs.upbit.com/docs/create-authorization-request

func CreateTokenWithNoParams() string {

	payload := jwt.MapClaims{
		"access_key": config.AccessKey,
		"nonce":      uuid.New().String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		fmt.Println("Error generating JWT token:", err)
		panic(err)
	}

	// Construct the authorization token
	authorizationToken := "Bearer " + tokenString

	return authorizationToken
}

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
		fmt.Println("Error generating JWT token:", err)
		panic(err)
	}

	// Construct the authorization token
	authorizationToken := "Bearer " + tokenString

	return authorizationToken
}

// OrderInfoAuthToken 주문가능 정보
// Deprecated
// https://docs.upbit.com/reference/%EC%A3%BC%EB%AC%B8-%EA%B0%80%EB%8A%A5-%EC%A0%95%EB%B3%B4
func OrderInfoAuthToken() {
	body := url.Values{}
	body.Set("market", "KRW-BTC")
	query := body.Encode()

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
		fmt.Println("Error generating JWT token:", err)
		panic(err)
	}

	// Construct the authorization token
	authorizationToken := "Bearer " + tokenString

	client := &http.Client{}
	serverURL := fmt.Sprintf("https://api.upbit.com/v1/orders/chance?%s", query)
	req, err := http.NewRequest("GET", serverURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		panic(err)
	}
	req.Header.Add("Authorization", authorizationToken)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		panic(err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(respBody))
}
