package users

import (
	"log"

	"github.com/enwawerueli/fiber-api/database"
	"github.com/enwawerueli/fiber-api/models"
	"github.com/enwawerueli/fiber-api/utils"
	"github.com/enwawerueli/fiber-api/validator"
	"github.com/gofiber/fiber/v2"
)

type UserIn struct {
	Username        string `json:"username" validate:"required" extensions:"x-order=0"`
	Email           string `json:"email" validate:"required,email" extensions:"x-order=1"`
	Password        string `json:"password" validate:"required" extensions:"x-order=2"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=password" copy:"-" extensions:"x-order=3"`
}

type UserOut struct {
	ID       uint   `json:"id" extensions:"x-order=0"`
	Username string `json:"username" extensions:"x-order=1"`
	Email    string `json:"email" extensions:"x-order=2"`
}

func (u *UserIn) PasswordHash() string {
	var hash, err = utils.HashPassword(u.Password)
	if err != nil {
		log.Fatalf("An error occurred while encryting password.\n%s", err.Error())
	}
	return hash
}

type Query struct {
	Page       int
	Size       int
	Properties []string
}

// @summary      Get users
// @description  Get a list of users
// @tags         User
// @accept       json
// @produce      json
// @param        page        query    integer  false  "Page number"     default(1)
// @param        size        query    integer  false  "Items per page"  default(30)
// @param        properties  query    array    false  "Item properties to fetch"
// @success      200         {array}  UserOut
// @security     ApiKeyAuth
// @router       /api/users [get]
func GetAll(c *fiber.Ctx) error {
	var q Query
	c.QueryParser(&q)
	var users = []models.User{}
	var offset int
	if q.Page > 1 {
		offset = (q.Page - 1) * q.Size
	}
	database.DB.Select(q.Properties).Limit(q.Size).Offset(offset).Find(&users)
	var usersOut = []UserOut{}
	for _, user := range users {
		var userOut UserOut
		utils.Copy(&user, &userOut)
		usersOut = append(usersOut, userOut)
	}
	return c.JSON(usersOut)
}

// @summary      Get user
// @description  Get a single user
// @tags         User
// @accept       json
// @produce      json
// @param        id   path      integer  true  "ID"
// @success      200  {object}  UserOut
// @security     ApiKeyAuth
// @router       /api/users/{id} [get]
func GetOne(c *fiber.Ctx) error {
	var id, err = c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.ErrBadRequest)
	}
	var user models.User
	database.DB.Find(&user, id)
	if user.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.ErrNotFound)
	}
	var userOut UserOut
	utils.Copy(&user, &userOut)
	return c.JSON(userOut)
}

// @summary      Create user
// @description  Create a new user
// @tags         User
// @accept       json
// @produce      json
// @param        user  body      UserIn  true  "New User"
// @success      200   {object}  UserOut
// @security     ApiKeyAuth
// @router       /api/users [post]
func Create(c *fiber.Ctx) error {
	var userIn UserIn
	if err := c.BodyParser(&userIn); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.ErrBadRequest)
	}
	if errs := validator.ValidateStruct(&userIn); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}
	var user models.User
	utils.Copy(&userIn, &user)
	database.DB.Create(&user)
	var userOut UserOut
	utils.Copy(&user, &userOut)
	return c.Status(201).JSON(userOut)
}
