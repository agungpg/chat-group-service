package auth

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

type MiddleWare struct {
}

func NewMiddleware() *MiddleWare {
	return &MiddleWare{}
}
func (m *MiddleWare) JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		fmt.Println("JWTMiddleware is running")
		auth := c.Get("Authorization")
		fmt.Println("auth", auth)
		if len(auth) < 8 || auth[:7] != "Bearer " {
			return c.Status(401).JSON(fiber.Map{"error": "missing or invalid token"})
		}
		tokenStr := auth[7:]

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "invalid token"})
		}

		// Store claims in context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Locals("userClaims", claims)
		}
		fmt.Println("pass token validation")
		// You can add `c.Locals("user_id", userId)` if needed
		return c.Next()
	}
}
