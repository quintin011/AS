package encryption

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/cw2/backend/models"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/matelang/jwt-go-aws-kms/v2/jwtkms"
)

func GenToken(User *models.User) *jwt.Token {
	t := jwt.NewWithClaims(jwtkms.SigningMethodRS256, &jwt.RegisteredClaims{
		Issuer: "Prod",
		Subject: User.UID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
	})
	return t
}

func Signstring(t *jwt.Token) string {

	s, err := t.SignedString(kmscfg.WithContext(context.Background()))
	if err != nil {
		log.Panic(err)
	}
	return s
}

func ParseToken(s string) *jwt.Token {
	t, err := jwt.Parse(s, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return kmscfg, nil
	})
	if err != nil {
		//log.Panic(err)
	}
	return t
}

func RefreshToken(t *jwt.Token) string {
	claims, ok := t.Claims.(jwt.MapClaims)
	if ok {
		claims["exp"] = jwt.NewNumericDate(time.Now().Add(time.Minute * 30))
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	return Signstring(token)
}

func SplitJWT(s string) string {
	jwtToken := strings.Split(s, " ")
	if len(jwtToken) != 2 {
		return ""
	}
	return jwtToken[1]
}
