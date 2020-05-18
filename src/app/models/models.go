package models

import (
	"github.com/jinzhu/gorm"
	"sort"
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
	ID        uint `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	User        User `gorm:"foreignkey:Userid" json:"user"`
	Userid      int
	Post	Post	`gorm:"foreignkey:Postid"`
	Postid	int
	Content     string `gorm:"type:TEXT;not null" json:"content"`
}

type ProfileUpdate struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword string `json:"newPassword"`
	Email	string	`json:"email"`
	FirstName	string	`json:"first_name"`
	LastName	string	`json:"last_name"`
	Password	string	`json:"password"`
}

type Chat struct {
	sort.Interface
	ID        uint `gorm:"primary_key" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	Content	string	`gorm:"not null" json:"content"`
	Sender	string	`gorm:"not null" json:"sender"`
	Receiver	string	`gorm:"not null" json:"receiver"`
}

type ChatsById []Chat

func (a ChatsById) Less(i, j int) bool {
	return a[i].ID < a[j].ID
}

func (a ChatsById) Len() int {
	return len(a)
}

func (a ChatsById) Swap(i, j int)  {
	a[i], a[j] = a[j], a[i]
}

func Migrate()  {

	Db.AutoMigrate(&User{},&Post{},&OauthClientDetails{},&Comment{},&Chat{})
	Db.Model(&Post{}).AddForeignKey("userid", "users", "RESTRICT", "RESTRICT")
	Db.Model(&Comment{}).AddForeignKey("userid", "users", "RESTRICT", "RESTRICT")
	Db.Model(&Comment{}).AddForeignKey("postid", "posts", "RESTRICT", "RESTRICT")
}

type SocketInfo struct {
	Username string	`json:"username"`
	FirstName	string `json:"first_name"`
	LastName	string `json:"last_name"`
	Name	string	`json:"name"`
	ID		string
}


type ChatMessage struct {
	Message string `json:"message" binding:"required"`
	To string `json:"to" binding:"required"`
}

type CommentModel struct {
	Content string `json:"content"	binding:"required"`
	PostId  int    `json:"postId"	binding:"required"`
}