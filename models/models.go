package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username     string `gorm:"not null;unique"`
	Email        string `gorm:"not null"`
	PasswordHash string `gorm:"not null"`
	Posts        []*Post
	Comments     []*Comment
}

type Post struct {
	gorm.Model
	Title    string `gorm:"not null;unique"`
	Content  string `gorm:""`
	UserID   uint   `gorm:"not null"`
	Author   *User  `gorm:"foreignKey:UserID"`
	Comments []*Comment
}

type Comment struct {
	gorm.Model
	Content string
	PostID  uint `gorm:"not null"`
	Post    *Post
	UserID  uint  `gorm:"not null"`
	Author  *User `gorm:"foreignKey:UserID"`
}
