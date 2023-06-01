package rs256

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	UserID     int
	Username   string
	GrantScope string
}

type MyCustomClaims struct {
	User
	jwt.RegisteredClaims
}

// 随机字符串
var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStr(str_len int) string {
	rand_bytes := make([]rune, str_len)
	for i := range rand_bytes {
		rand_bytes[i] = letters[rand.Intn(len(letters))]
	}
	return string(rand_bytes)
}

// pkcs1
func parsePriKeyBytes(buf []byte) (*rsa.PrivateKey, error) {
	p := &pem.Block{}
	p, buf = pem.Decode(buf)
	if p == nil {
		return nil, errors.New("parse key error")
	}
	return x509.ParsePKCS1PrivateKey(p.Bytes)
}

const pri_key = `-----BEGIN RSA PRIVATE KEY-----
MIIBOwIBAAJBAJS/Q6xhgqzycc5KEPtQrI8Bf8WakGmpum3zutno28qASP1Y0BAs
WTEVNtrxuVXb9R3ebcrI5lYnkxR5DobffWkCAwEAAQJAdGptMJDwkSL+5xEY0ViG
dTYbJjCeLdRk0IEdEEcrHgSvn6YRvlQ5lsSMSFxRscCO40XMq9OofPzV6RD1kWgQ
AQIhAMXwMHC3qXmCcsxcyJ3yjCjRzSmA+PQLkijX2Z4nFQ/pAiEAwGEohnQXgdDc
vGauBJQc3clb33s4y+Woi86WEAkUEYECIHxCWsaIJfZH9DVjEfZF68M8YiVp99+M
3AaT6uOj+U7xAiEAk6n/8TQq1vn6dKJb8CfAAH0Oh/uNHPSq6qUniidtwAECIQCW
cPfmrnitlRFkmV9/c5Waln9VXpj0e0sU7Zj4GqOOSg==
-----END RSA PRIVATE KEY-----
`

const pub_key = `-----BEGIN RSA PUBLIC KEY-----
MEgCQQCUv0OsYYKs8nHOShD7UKyPAX/FmpBpqbpt87rZ6NvKgEj9WNAQLFkxFTba
8blV2/Ud3m3KyOZWJ5MUeQ6G331pAgMBAAE=
-----END RSA PUBLIC KEY-----
`

func GenerateTokenUsingRS256(user User) (string, error) {
	claim := MyCustomClaims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Auth_Server",                                   // 签发者
			Subject:   "Tom",                                           // 签发对象
			Audience:  jwt.ClaimStrings{"Android_APP", "IOS_APP"},      //签发受众
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),   //过期时间
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second)), //最早使用时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                  //签发时间
			ID:        randStr(10),                                     // jwt ID, 类似于盐值
		},
	}
	rsa_pri_key, err := parsePriKeyBytes([]byte(pri_key))
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claim).SignedString(rsa_pri_key)
	return token, err
}

func parsePubKeyBytes(pub_key []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pub_key)
	if block == nil {
		return nil, errors.New("block nil")
	}
	pub_ret, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, errors.New("x509.ParsePKCS1PublicKey error")
	}

	return pub_ret, nil
}

func ParseTokenRs256(token_string string) (*MyCustomClaims, error) {
	token, err := jwt.ParseWithClaims(token_string, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		pub, err := parsePubKeyBytes([]byte(pub_key))
		if err != nil {
			fmt.Println("err = ", err)
			return nil, err
		}
		return pub, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("claim invalid")
	}

	claims, ok := token.Claims.(*MyCustomClaims)
	if !ok {
		return nil, errors.New("invalid claim type")
	}

	return claims, nil
}
