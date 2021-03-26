package common

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"oceanlearn.teach/ginessential/model"
)

// 定义加密秘钥
var jwtKey = []byte("a_secret_crect")

// 定义token的claim
type Claims struct {
	UserId uint
	jwt.StandardClaims
}

// 定义方法,登录成功后调用这个方法发放token
func ReleaseToken(user model.User) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour) // token的过期时间，这里是7天
	claims := &Claims{
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), // 过期时间
			IssuedAt:  time.Now().Unix(),     // 发放时间
			Issuer:    "oceanlearn.tech",     // 发放者
			Subject:   "user token",          // 主题
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})
	return token, claims, err
}
