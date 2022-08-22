package auth

import (
	"github.com/enwawerueli/fiber-api/database"
	"github.com/enwawerueli/fiber-api/middleware/session"
	"github.com/enwawerueli/fiber-api/models"
	"github.com/enwawerueli/fiber-api/utils"
	"github.com/gofiber/fiber/v2"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *fiber.Ctx) error {
	var userIn User
	if err := c.BodyParser(&userIn); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.ErrBadRequest)
	}
	var user models.User
	database.DB.Find(&user, "username=?", userIn.Username)
	if user.ID == 0 || !utils.VerifyPassword(user.PasswordHash, userIn.Password) {
		return c.JSON("Invalid username or password")
	}
	sess, err := session.Store.Get(c)
	if err != nil {
		panic(err)
	}
	sess.Set("uid", user.Username)
	if err := sess.Save(); err != nil {
		panic(err)
	}
	return c.JSON("Login successful")
}

func Logout(c *fiber.Ctx) error {
	sess, err := session.Store.Get(c)
	if err != nil {
		panic(err)
	}
	// Destroy session
	if err := sess.Destroy(); err != nil {
		panic(err)
	}
	return c.JSON("Logged out")
}

func Authorize(c *fiber.Ctx) error {
	return nil
}
