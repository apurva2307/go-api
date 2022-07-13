package middlewares

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type TokenClaims struct {
	Name string
	Exp  int64
	jwt.MapClaims
}

func AuthenticateToken(c *fiber.Ctx) error {
	token := c.GetReqHeaders()["Token"]

	if token == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors": map[string]string{"msg": "You are not authorized to access this route."},
		})
	}
	decodedToken, err := jwt.ParseWithClaims(token, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors": map[string]string{"msg": "Authorization token is invalid."},
		})
	}
	claims, _ := decodedToken.Claims.(*TokenClaims)
	timeNow := time.Now().Unix()
	if timeNow > claims.Exp || !decodedToken.Valid {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors": map[string]string{"msg": "Authorization token is expired or invalid."},
		})
	}
	if claims.Name != "shailendra" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors": map[string]string{"msg": "Authorization token is invalid."},
		})
	}
	return c.Next()
}
