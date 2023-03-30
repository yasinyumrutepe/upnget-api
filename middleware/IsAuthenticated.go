package middleware

import (
	"auction/controllers"
	"auction/globals"
	"auction/models"
	"auction/secret"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func IsManuelAuthenticated(c *fiber.Ctx) error {
	JWT := secret.Env["JWT"].(map[string]interface{})
	//DESC - JWT control
	tokenString, ok := c.GetReqHeaders()["Authorization"]
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "JWT missing!!!", "data": nil})
	}
	tknsplit := strings.Split(tokenString, " ")
	tknstr, err := changeVar(tknsplit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Header parsing error!!!", "data": nil})
	}
	token, err := jwt.Parse(tknstr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte([]byte(JWT["secret"].(string))), nil
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": err.Error(), "data": nil})
	}

	//DESC - JWT control END
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["user_id"] == "" {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "User not found", "data": nil})
		}
		//DESC - get user in jwt for check role
		userAbility := models.User{
			ID:    uint(claims["user_id"].(float64)),
			Level: uint(claims["level"].(float64)),
		}
		c.Locals("user", userAbility)

		return c.Next()
	} else {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": err.Error(), "data": nil})
	}
}

func IsAutoAuthenticated(c *fiber.Ctx) error {

	JWT := secret.Env["JWT"].(map[string]interface{})
	//DESC - JWT control
	tokenString, ok := c.GetReqHeaders()["Authorization"]
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "JWT missing!!!", "data": nil})
	}
	tknsplit := strings.Split(tokenString, " ")
	tknstr, err := changeVar(tknsplit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Header parsing error!!!", "data": nil})
	}
	token, err := jwt.Parse(tknstr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte([]byte(JWT["secret"].(string))), nil
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": err.Error(), "data": nil})
	}

	//DESC - JWT control END
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["user_id"] == "" {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "User not found", "data": nil})
		}
		//DESC - get user in jwt for check role
		userAbility := models.User{
			ID:    uint(claims["user_id"].(float64)),
			Level: uint(claims["level"].(float64)),
		}
		c.Locals("user", userAbility)

		//DESC - User Ability Control
		//DESC - dynamic method ability control
		method := c.Method()
		path := c.Path()
		res1 := strings.Split(path, "/")
		//DESC - person ability control function
		control, err := controllers.HasAbility(method, res1[2], c)

		if err != nil || !control {
			return c.Status(418).JSON(globals.Response{
				Error:   true,
				Message: " The server refuses to brew coffee because it is, permanently, a teapot.",
			})
		}

		return c.Next()
	} else {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": err.Error(), "data": nil})
	}

}

func changeVar(arr []string) (string, error) {
	var s string
	if len(arr) == 1 {
		s = arr[0]
		return s, nil
	} else if len(arr) == 2 {
		s = arr[1]
		return s, nil
	} else {
		return "", fmt.Errorf("Split error")
	}
}
