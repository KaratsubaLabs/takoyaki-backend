package main

import (
    "os"
    "errors"
    "time"
    "github.com/golang-jwt/jwt"
)

var JWTLifetime = 4096
var JWTSecret = os.Getenv("JWT_SECRET")

type Claims struct {
    UserID uint
    jwt.StandardClaims
}

func GenerateToken(id uint) (string, error) {

    expire := time.Now().Add(time.Duration(JWTLifetime) * time.Minute)
    claims := &Claims{
        UserID: id,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expire.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(JWTSecret))
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

func ValidateToken(tokenString string) (uint, error) {

    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return []byte(JWTSecret), nil
    })
    if err != nil {
        return 0, err
    }
    if !token.Valid {
        return 0, errors.New("invalid access token")
    }
    return claims.UserID, nil
}

