package common

import (
	"clock/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// jwt加密密钥
var jwtKey = []byte("liudebaoliliuqin")

type Claims struct {
	UID string
	jwt.RegisteredClaims
}

// 生成Token
func ReleaseToken(user model.User) (string, error) {
	//token的有效期
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		//自定义字段
		UID: user.UserID,
		//标准字段
		RegisteredClaims: jwt.RegisteredClaims{
			//过期时间
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			//发放时间
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}
	//使用jwt密钥生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	//返回token
	return tokenString, nil
}

// 解析token
func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, nil, err
	}

	// 检查 token 是否过期
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, nil, jwt.ErrTokenExpired
	}

	return token, claims, nil
}
