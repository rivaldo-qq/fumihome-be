package jwt

import (
	"context"
	"fmt"
	"os"

	"github.com/Dryluigi/go-grpc-ecommerce-be/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

type JwtEntityContextKey string

var JwtEntityContextKeyValue JwtEntityContextKey = "JwtEntity"

type JwtClaims struct {
	jwt.RegisteredClaims
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

func (jc *JwtClaims) SetToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, JwtEntityContextKeyValue, jc)
}

func GetClaimsFromToken(token string) (*JwtClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &JwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, utils.UnauthenticatedResponse()
	}

	if !tokenClaims.Valid {
		return nil, utils.UnauthenticatedResponse()
	}

	if claims, ok := tokenClaims.Claims.(*JwtClaims); ok {
		return claims, nil
	}

	return nil, utils.UnauthenticatedResponse()
}

func GetClaimsFromContext(ctx context.Context) (*JwtClaims, error) {
	claims, ok := ctx.Value(JwtEntityContextKeyValue).(*JwtClaims)
	if !ok {
		return nil, utils.UnauthenticatedResponse()
	}

	return claims, nil
}
