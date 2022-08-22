package session

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/sqlite3"
)

var (
	Store *session.Store
)

func init() {
	var storage = sqlite3.New(sqlite3.Config{
		Database: "./db.sqlite",
		Table:    "sessions",
	})
	Store = session.New(session.Config{
		Storage:    storage,
		Expiration: 30 * time.Minute,
	})
}

func IsAuthenticated(c *fiber.Ctx) error {
	sess, err := Store.Get(c)
	if err != nil {
		panic(err)
	}
	if sess.Get("uid") == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.ErrUnauthorized)
	}
	return c.Next()
}
