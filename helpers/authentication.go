package helpers

import (
	"os"
	"time"

	"github.com/apurva2307/go-api/models"
	"github.com/golang-jwt/jwt/v4"
)

func CreateJwtToken(user models.User) (string, error) {
	exp := time.Now().Add(300 * time.Second).Unix()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = exp
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return t, err
}
