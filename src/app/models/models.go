package models

import (
	"github.com/jinzhu/gorm"
	"time"
)
import 	_ "github.com/jinzhu/gorm/dialects/postgres"

var Db = initializeDB()

func initializeDB() *gorm.DB  {
	db, err := gorm.Open("postgres", "host=127.0.0.1 port=5432 user=postgres dbname=go_social password=root sslmode=disable")
	if err != nil {
		panic( err.Error())
	}

	return db
}

type OauthClientDetails struct {
	gorm.Model
	ClientId     string `gorm:"unique;not null"`
	ClientSecret string `gorm:"unique;not null"`
	Name         string `gorm:"unique;not null"`
}

type User struct {
	gorm.Model
	Username       string `gorm:"type:VARCHAR(20);unique;not null" json:"username" binding:"required"`
	FirstName      string `gorm:"type:VARCHAR(20);not null"  json:"first_name" binding:"required"`
	LastName       string `gorm:"type:VARCHAR(20);not null"  json:"last_name" binding:"required"`
	Password       string `gorm:"type:VARCHAR(255);not null"  json:"-"`
	Gender         string `gorm:"type:VARCHAR(10);not null"  json:"gender" binding:"required"`
	Email          string `gorm:"type:VARCHAR(100);not null" json:"email" binding:"required"`
	DisplayPicture string `json:"display_picture"`
}

type Post struct {
	ID        uint `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	User        User `gorm:"foreignkey:Userid" json:"user"`
	Userid      int
	Content     string `gorm:"type:TEXT;not null" json:"content"`
	Attachments string `gorm:"type:TEXT;" json:"attachments"`
}

type Comment struct {
	gorm.Model
	User        User `gorm:"foreignkey:Userid"`
	Userid      int
	Post	Post	`gorm:"foreignkey:Postid"`
	Postid	int
	Content     string `gorm:"type:TEXT;not null"`
}

type ProfileUpdate struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword string `json:"newPassword"`
	Email	string	`json:"email"`
	FirstName	string	`json:"first_name"`
	LastName	string	`json:"last_name"`
	Password	string	`json:"password"`
}

func Migrate()  {

	Db.AutoMigrate(&User{},&Post{},&OauthClientDetails{},&Comment{})
	Db.Model(&Post{}).AddForeignKey("userid", "users", "RESTRICT", "RESTRICT")
	Db.Model(&Comment{}).AddForeignKey("userid", "users", "RESTRICT", "RESTRICT")
	Db.Model(&Comment{}).AddForeignKey("postid", "posts", "RESTRICT", "RESTRICT")
}

