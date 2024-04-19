package utils

import (
	"errors"
	"fmt"
	"github.com/go-chi/jwtauth/v5"
	"github.com/golang-jwt/jwt/v5"
	"kiramishima/ionix/internal/models"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var privateKey = []byte(os.Getenv("JWT_PRIVATE_KEY"))
var TokenAuth = jwtauth.New("HS256", []byte(os.Getenv("JWT_PRIVATE_KEY")), nil)

// GenerateJWT generate JWT token
func GenerateJWT(user *models.User) (string, error) {
	tokenTTL, _ := strconv.Atoi(os.Getenv("TOKEN_TTL"))
	// fmt.Println("TokenTTL: ", tokenTTL)
	// fmt.Println("JWT_PRIVATE_KEY: ", os.Getenv("JWT_PRIVATE_KEY"))
	userID := strconv.Itoa(int(user.ID))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		// A usual scenario is to set the expiration time relative to the current time
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(tokenTTL) * time.Second)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		ID:        userID,
	})
	return token.SignedString(privateKey)
}

// ValidateJWT validate JWT token
func ValidateJWT(req *http.Request) error {
	token, err := getToken(req)
	if err != nil {
		return err
	}
	_, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return nil
	}
	return errors.New("invalid token provided")
}

// ValidateAdminRoleJWT validate Admin role
func ValidateAdminRoleJWT(req *http.Request) error {
	token, err := getToken(req)
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	userRole := uint(claims["role"].(float64))
	if ok && token.Valid && userRole == 1 {
		return nil
	}
	return errors.New("invalid admin token provided")
}

// ValidateCustomerRoleJWT validate Customer role
func ValidateCustomerRoleJWT(req *http.Request) error {
	token, err := getToken(req)
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	userRole := uint(claims["role"].(float64))
	if ok && token.Valid && userRole == 2 || userRole == 1 {
		return nil
	}
	return errors.New("invalid customer or admin token provided")
}

func GetUserIDInJWTHeader(req *http.Request) int {
	_, decoded, _ := jwtauth.FromContext(req.Context())
	log.Println(decoded)
	var sID = decoded["id"]
	var ID, _ = strconv.ParseInt(sID.(string), 10, 32)
	return int(ID)
}

// getToken check token validity
func getToken(req *http.Request) (*jwt.Token, error) {
	tokenString := getTokenFromRequest(req)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return privateKey, nil
	})
	return token, err
}

// getTokenFromRequest extract token from request Authorization header
func getTokenFromRequest(req *http.Request) string {
	bearerToken := req.Header.Get("Authorization")
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}
