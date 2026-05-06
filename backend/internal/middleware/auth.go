package middleware

import (
	"strings"

	"github.com/fitrianabila2025group/videoxnx/backend/internal/config"
	"github.com/fitrianabila2025group/videoxnx/backend/internal/utils"
	"github.com/gofiber/fiber/v2"
)

// RequireAdmin verifies a JWT token from Authorization: Bearer or cookie.
func RequireAdmin(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := ""
		auth := c.Get("Authorization")
		if strings.HasPrefix(auth, "Bearer ") {
			token = strings.TrimPrefix(auth, "Bearer ")
		}
		if token == "" {
			token = c.Cookies("admin_token")
		}
		if token == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "missing token")
		}
		claims, err := utils.ParseJWT(cfg.JWTSecret, token)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
		}
		if claims.Role != "admin" {
			return fiber.NewError(fiber.StatusForbidden, "forbidden")
		}
		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)
		return c.Next()
	}
}

// SecureHeaders sets a baseline of security headers.
func SecureHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "SAMEORIGIN")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		return c.Next()
	}
}
