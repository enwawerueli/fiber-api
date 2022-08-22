package main

import (
	"log"

	_ "github.com/enwawerueli/fiber-api/docs"
	"github.com/enwawerueli/fiber-api/handlers/auth"
	"github.com/enwawerueli/fiber-api/handlers/comments"
	"github.com/enwawerueli/fiber-api/handlers/posts"
	"github.com/enwawerueli/fiber-api/handlers/users"
	"github.com/enwawerueli/fiber-api/middleware/session"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func root(c *fiber.Ctx) error {
	return c.JSON("Welcome to My Fiber API")
}

// @title                       My Fiber API
// @version                     1.0
// @description                 Fiber API with auto-generated  Swagger docs.
// @termsOfService              http://swagger.io/terms/
// @contact.name                Support
// @contact.email               seaworndrift@gmail.com
// @license.name                Apache 2.0
// @license.url                 http://www.apache.org/licenses/LICENSE-2.0.html
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
// @basePath                    /api
// @schemes                     http
func main() {
	// Initialize app
	var app = fiber.New()

	// Set up routes
	app.Get("/", session.IsAuthenticated, root)
	// Swagger docs
	app.Get("/docs/*", swagger.HandlerDefault)
	var api = app.Group("/api")
	// Auth
	app.Post("/login", auth.Login)
	app.Get("/logout", auth.Logout)
	app.Post("/register", users.Create)
	// Users
	api.Get("/users", users.GetAll)
	api.Get("/users/:id", users.GetOne)
	api.Post("/users", users.Create)
	// Posts
	api.Get("/posts", posts.GetAll)
	api.Get("/posts/:id", posts.GetOne)
	api.Post("/posts", posts.Create)
	// Comments
	api.Get("/comments", comments.GetAll)
	api.Get("/comments/:id", comments.GetOne)
	api.Post("/comments", comments.Create)

	// Run server
	log.Fatal(app.Listen(":3000"))
}
