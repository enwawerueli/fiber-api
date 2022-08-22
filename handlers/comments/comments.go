package comments

import (
	"github.com/enwawerueli/fiber-api/database"
	"github.com/enwawerueli/fiber-api/models"
	"github.com/enwawerueli/fiber-api/utils"
	"github.com/enwawerueli/fiber-api/validator"
	"github.com/gofiber/fiber/v2"
)

type CommentIn struct {
	Content string `json:"content" validate:"required" extensions:"x-order=0"`
	PostID  uint   `json:"post_id" validate:"required,min=1" extensions:"x-order=1"`
	UserID  uint   `json:"user_id" validate:"required,min=1" extensions:"x-order=2"`
}

type CommentOut struct {
	ID      uint   `json:"id" extensions:"x-order=0"`
	Content string `json:"content" extensions:"x-order=1"`
	PostID  uint   `json:"post_id" extensions:"x-order=2"`
	UserID  uint   `json:"user_id" extensions:"x-order=3"`
}

type Query struct {
	Page       int
	Size       int
	Properties []string
}

// @summary      Get comments
// @description  Get a list of comments
// @tags         Comment
// @accept       json
// @produce      json
// @param        page        query    integer  false  "Page number"     default(1)
// @param        size        query    integer  false  "Items per page"  default(30)
// @param        properties  query    array    false  "Item properties to fetch"
// @success      200         {array}  CommentOut
// @security     ApiKeyAuth
// @router       /api/comments [get]
func GetAll(c *fiber.Ctx) error {
	var q Query
	c.QueryParser(&q)
	var comments = []models.Comment{}
	var offset int
	if q.Page > 1 {
		offset = (q.Page - 1) * q.Size
	}
	database.DB.Select(q.Properties).Limit(q.Size).Offset(offset).Find(&comments)
	var commentsOut = []CommentOut{}
	for _, comment := range comments {
		var commentOut CommentOut
		utils.Copy(&comment, &commentOut)
		commentsOut = append(commentsOut, commentOut)
	}
	return c.JSON(commentsOut)
}

// @summary      Get comment
// @description  Get a single comment
// @tags         Comment
// @accept       json
// @produce      json
// @param        id   path      integer  true  "ID"
// @success      200  {object}  CommentOut
// @security     ApiKeyAuth
// @router       /api/comments/{id} [get]
func GetOne(c *fiber.Ctx) error {
	var id, err = c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.ErrBadRequest)
	}
	var comment models.Comment
	database.DB.Find(&comment, id)
	if comment.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.ErrNotFound)
	}
	var commentOut CommentOut
	utils.Copy(&comment, &commentOut)
	return c.JSON(commentOut)
}

// @summary      Create comment
// @description  Create a new comment
// @tags         Comment
// @accept       json
// @produce      json
// @param        comment  body      CommentIn  true  "New Comment"
// @success      200      {object}  CommentOut
// @security     ApiKeyAuth
// @router       /api/comments [post]
func Create(c *fiber.Ctx) error {
	var commentIn CommentIn
	if err := c.BodyParser(&commentIn); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.ErrBadRequest)
	}
	if errs := validator.ValidateStruct(&commentIn); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}
	var comment models.Comment
	utils.Copy(&commentIn, &comment)
	database.DB.Create(&comment)
	var commentOut CommentOut
	utils.Copy(&comment, &commentOut)
	return c.Status(201).JSON(commentOut)
}
