package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gfmanica/splitz-backend/config"
	"github.com/gfmanica/splitz-backend/types"
	"github.com/gfmanica/splitz-backend/utils"
	"github.com/golang-jwt/jwt"
)

type contextKey string

const UserKey contextKey = "id"

func CreateJWT(secret []byte, user *types.User) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        strconv.Itoa(user.ID),
		"name":      user.Name,
		"email":     user.Email,
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := getTokenFromRequest(r)

		token, err := validateToken(tokenString)

		if err != nil {
			log.Printf("Error validating token: %v", err)

			permissionDenied(w)

			return
		}

		if !token.Valid {
			log.Printf("invalid token")

			permissionDenied(w)

			return
		}

		claims := token.Claims.(jwt.MapClaims)
		str := claims["id"].(string)

		userID, _ := strconv.Atoi(str)

		u, err := store.GetUserByID(userID)

		if err != nil {
			log.Printf("Error getting user by ID: %v", err)

			permissionDenied(w)

			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, u)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}

func getTokenFromRequest(r *http.Request) string {
	token := r.Header.Get("Authorization")

	if token != "" {
		return token
	}

	return ""
}

func validateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.Envs.JWTSecret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriterError(w, http.StatusUnauthorized, fmt.Errorf("permission denied"))
}

func GetUserIDFromContext(ctx context.Context) int {
	userID, ok := ctx.Value(UserKey).(int)

	if !ok {
		return -1
	}

	return userID
}
