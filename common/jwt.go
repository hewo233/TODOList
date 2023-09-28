package common

import (
	"TODOList/model"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtKey = []byte("this_is_hewoKEY")

type Claim struct {
	UserId uint
	jwt.StandardClaims
}

func ReleaseToken(user model.User) (string, error) {
	//token 的有效期
	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	claims := &Claim{
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), //过期时间
			IssuedAt:  time.Now().Unix(),     //发放时间
			Issuer:    "127.0.0.1",
			Subject:   "user token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, *Claim, error) {
	claim := &Claim{}
	token, err := jwt.ParseWithClaims(tokenString, claim, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})
	return token, claim, err
}
