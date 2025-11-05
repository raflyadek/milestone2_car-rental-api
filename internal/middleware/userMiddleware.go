package middleware

import (
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func JwtMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
		jwtMiddleware := echojwt.WithConfig(echojwt.Config{
			SigningKey: []byte(jwtSecretKey),
		})
		return jwtMiddleware(next)(c)
	}
}