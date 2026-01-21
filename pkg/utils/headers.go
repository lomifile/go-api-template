package utils

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ExtractJwtTokenFromHeader takes fiber.Ctx and extract Authorization header from it
// and returns token portion of JWT
func ExtractJwtTokenFromHeader(c *fiber.Ctx) (string, error) {
	headers := c.GetReqHeaders()

	authorization, ok := headers["Authorization"]

	if !ok {
		return "", errors.New("headers doesn't contain token")
	}

	fullToken := authorization[0]

	split := strings.Split(fullToken, " ")
	token := split[1]

	return token, nil
}
