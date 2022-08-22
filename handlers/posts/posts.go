package posts

import (
	"github.com/enwawerueli/fiber-api/database"
	"github.com/enwawerueli/fiber-api/models"
	"github.com/enwawerueli/fiber-api/utils"
	"github.com/enwawerueli/fiber-api/validator"
	"github.com/gofiber/fiber/v2"
)

type PostIn struct {
	Title   string `json:"title" validate:"required" extensions:"x-order=0"`
	Content string `json:"content" validate:"required" extensions:"x-order=1"`
	UserID  string `json:"user_id" validate:"required,min=1" extensions:"x-order=2"`
}

type PostOut struct {
	ID      uint   `json:"id" extensions:"x-order=0"`
	Title   string `json:"title" extensions:"x-order=1"`
	Content string `json:"content" extensions:"x-order=2"`
	UserID  uint   `json:"user_id" extensions:"x-order=3"`
}

type Query struct {
	Page       int
	Size       int
	Properties []string
}

// @summary      Get posts
// @description  Get a list of posts
// @tags         Post
// @accept       json
// @produce      json
// @param        page        query    integer  false  "Page number"     default(1)
// @param        size        query    integer  false  "Items per page"  default(30)
// @param        properties  query    array    false  "Item properties to fetch"
// @success      200         {array}  PostOut
// @security     ApiKeyAuth
// @router       /api/posts [get]
func GetAll(c *fiber.Ctx) error {
	var q Query
	c.QueryParser(&q)
	var posts = []models.Post{}
	var offset int
	if q.Page > 1 {
		offset = (q.Page - 1) * q.Size
	}
	database.DB.Select(q.Properties).Limit(q.Size).Offset(offset).Find(&posts)
	var postsOut = []PostOut{}
	for _, post := range posts {
		var postOut PostOut
		utils.Copy(&post, &postOut)
		postsOut = append(postsOut, postOut)
	}
	return c.JSON(postsOut)
}

// @summary      Get post
// @description  Get a single post
// @tags         Post
// @accept       json
// @produce      json
// @param        id   path      integer  true  "ID"
// @success      200  {object}  PostOut
// @security     ApiKeyAuth
// @router       /api/posts/{id} [get]
func GetOne(c *fiber.Ctx) error {
	var id, err = c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.ErrBadRequest)
	}
	var post models.Post
	database.DB.Find(&post, id)
	if post.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.ErrNotFound)
	}
	var postOut PostOut
	utils.Copy(&post, &postOut)
	return c.JSON(postOut)
}

// @summary      Create post
// @description  Create a new post
// @tags         Post
// @accept       json
// @produce      json
// @Param        post  body      PostIn  true  "Post"
// @success      200   {object}  PostOut
// @security     ApiKeyAuth
// @router       /api/posts [post]
func Create(c *fiber.Ctx) error {
	var postIn PostIn
	if err := c.BodyParser(&postIn); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.ErrBadRequest)
	}
	if errs := validator.ValidateStruct(&postIn); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}
	var post models.Post
	utils.Copy(&postIn, &post)
	database.DB.Create(&post)
	var postOut PostOut
	utils.Copy(&post, &postOut)
	return c.Status(201).JSON(postOut)
}
