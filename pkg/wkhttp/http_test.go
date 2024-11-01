package wkhttp

import (
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"testing"
)

// 后面改成从nacos获取
const privateKeyBase64Test = "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEbTgMAh0eBehpVw7P2OIksraGPk5YLsFC9YKMLSA/hnvFhInROEHbJNaKSJoyaWJT7Ehc/x80r6CRYum1ghVvhg=="
const privateKeyBase64Prod = "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEr0M9I/Huh43yGwodp40UJSRIktdfjNIfR9IcCO04/aMKcyKUJEQH7Xnk4msY6jux6wfiBPDRnMJ8OHruwGZ9+A=="

func TestVerifyAndGetUsername(t *testing.T) {
	tokenString := "eyJhbGciOiJFUzI1NiJ9.eyJzdWIiOiJlY2hvb29fbWFsbF9lcnAiLCJpc3MiOiJlY2hvb29fbWFsbCIsInVzZXJuYW1lIjoiaHVhbmd4aWFvY29uZyIsImF1dGhfdGltZSI6MTczMDI2MDEwMCwiZXhwIjoxNzMxMTI0MTAwLCJpYXQiOjE3MzAyNjAxMDB9.QiyR6Ns8o395Rwy8fysdoc0h6j1xcJyMHBiVifZe8gxUSy5Rmokblz6GNt7zah2f0gmPXHIgkvjE4HUUoqOOww"

	// 对Base64编码的密钥进行解码
	decodedKeyBytesTest, err := base64.StdEncoding.DecodeString(privateKeyBase64Test)
	decodedKeyBytesProd, err := base64.StdEncoding.DecodeString(privateKeyBase64Prod)
	if err != nil {
		fmt.Println("解码Base64密钥失败:", err)
		return
	}

	// 将解码后的字节数据解析为ECDSA私钥
	publicKeyTest, err := x509.ParsePKIXPublicKey(decodedKeyBytesTest)
	if err != nil {
		fmt.Println("解析私钥失败:", err)
		return
	}

	publicKeyProd, err := x509.ParsePKIXPublicKey(decodedKeyBytesProd)
	if err != nil {
		fmt.Println("解析私钥失败:", err)
		return
	}

	// 创建用于验证JWT Token的函数，传入私钥
	keyFuncTest := func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Method)
		}
		return publicKeyTest, nil
	}

	// 创建用于验证JWT Token的函数，传入私钥
	keyFuncProd := func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Method)
		}
		return publicKeyProd, nil
	}

	// 解析JWT Token
	token, err := jwt.Parse(tokenString, keyFuncTest)

	if err != nil {
		token, err = jwt.Parse(tokenString, keyFuncProd)
		if err != nil {
			fmt.Println("解析JWT Token失败:", err)
			return
		}
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("用户ID:", claims["sub"])
		fmt.Println("用户名:", claims["username"])
		fmt.Println("发行者:", claims["iss"])
		fmt.Println("授权时间:", claims["auth_time"])
		fmt.Println("过期时间:", claims["exp"])
		fmt.Println("签发时间:", claims["iat"])
	} else {
		fmt.Println("无效的JWT Token或无法获取声明信息")
	}
}
