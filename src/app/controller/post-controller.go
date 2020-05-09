package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-social-app/src/app/models"
	"go-social-app/src/app/s3"
	"go-social-app/src/app/services"
	"net/http"
	"strconv"
)
const S3BucketUrl = "https://social-app-bucket1.s3.amazonaws.com/social-app-images/"

func InitiatePostRoute(securedRoute *gin.RouterGroup) {
	securedRoute.POST("/post-with-attachment", func(context *gin.Context) {
		form, err := context.MultipartForm()
		files :=form.File["files"]

		content := context.PostForm("content")
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message" : "Unable to retrieve file"})
			return
		}
		var filePaths []string
		if len(files) > 0 {
			filePaths, err = s3.UploadFilesToS3(files)
		}

		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message" : err.Error()})
			return
		}
		var filePathString string
		if len(filePaths) > 0 {
			fileB, _ := json.Marshal(filePaths)
			filePathString = string(fileB)
		}

		username , _ := context.Get("oauth.credential")
		user, error := services.LoadUserByUsername(username.(string))
		if error != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message" : "No user found"})
			return
		}

		var post = models.Post{
			Userid:      int(user.ID),
			Content:     content,
			Attachments: filePathString,
			User: user,
		}

		postResult := models.Db.Create(&post)
		if postResult.Error != nil {
			fmt.Println(postResult.Error)
			context.JSON(http.StatusBadRequest, gin.H{"message":"Error occurred. please try again"})
			return
		}

		context.JSON(http.StatusOK, post)

	})

	securedRoute.GET("/posts/", func(context *gin.Context) {
		pageNo, pageErr := strconv.Atoi(context.Query("pageNo"))
		pageSize, pageSizeErr := strconv.Atoi(context.Query("pageSize"))

		if pageErr != nil {
			pageNo = 0
		}

		if pageSizeErr != nil {
			pageSize = 10
		}

		posts, err := services.LoadPostByPageNumberAndSize(pageNo, pageSize)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"message" : err.Error()})
			return
		}

		context.JSON(http.StatusOK, posts)

	})

	securedRoute.GET("/posts-with-lesser-id/", func(context *gin.Context) {
		id, idErr := strconv.Atoi(context.Query("id"))
		pageSize, pageSizeErr := strconv.Atoi(context.Query("pageSize"))

		if idErr != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"message" : idErr.Error()})
			return
		}

		if pageSizeErr != nil {
			pageSize = 10
		}

		posts, error := services.LoadPostsWithLesserId(id, pageSize)

		if error != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"message" : error.Error()})
			return
		}

		context.JSON(http.StatusOK, posts)

	})
}
