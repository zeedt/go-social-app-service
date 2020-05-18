package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-social-app/src/app/models"
	"go-social-app/src/app/services"
	"net/http"
	"strconv"
)

func InitiateCommentRoute(securedRoute *gin.RouterGroup) {

	securedRoute.POST("/comment", func(context *gin.Context) {
		var commentModel models.CommentModel
		err := context.BindJSON(&commentModel)
		if err != nil {
			fmt.Println("Error occurred due to ", err.Error())
			context.JSON(http.StatusBadRequest, gin.H{"message" : "Error occurred while adding comment due to " + err.Error()})
			return
		}

		username , _ := context.Get("oauth.credential")
		user, error := services.LoadUserByUsername(username.(string))
		if error != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message" : "No user found"})
			return
		}

		var comment = models.Comment{
			User:    user,
			Userid:  int(user.ID),
			Postid: commentModel.PostId,
			Content: commentModel.Content,
		}

		result := models.Db.Create(&comment)
		if result.Error != nil {
			fmt.Println(result.Error)
			context.JSON(http.StatusBadRequest, gin.H{"message" : "Unable to persist comment due to " + result.Error.Error()})
			return
		}

		context.JSON(http.StatusOK, gin.H{"successful" : true, "message" : "successfully added comment"})

	})

	securedRoute.GET("/comment/:postId", func(context *gin.Context) {
		postId, err := strconv.Atoi(context.Param("postId"))
		pageNo, pageErr := strconv.Atoi(context.Query("pageNo"))
		pageSize, pageSizeErr := strconv.Atoi(context.Query("pageSize"))

		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message" : "Invalid post id"})
			return
		}

		if pageErr != nil {
			pageNo = 0
		}

		if pageSizeErr != nil {
			pageSize = 10
		}

		comments, error := services.LoadCommentsByPostIdAndPageNoAndSize(postId, pageNo, pageSize)

		if error != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message" : "Error occurred due to " + error.Error()})
			return
		}
		context.JSON(http.StatusOK, comments)
	})

	securedRoute.GET("/comment/:postId/comments-with-lesser-id", func(context *gin.Context) {
		comentId, err := strconv.Atoi(context.Query("commentId"))
		postId, postIdErr := strconv.Atoi(context.Param("postId"))
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message" : "Invalid comment id"})
			return
		}
		if postIdErr != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message" : "Invalid post id"})
			return
		}
		comments, loadError := services.LoadCommentsWithLesserId(comentId, postId, 20)
		if loadError != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"message":"Error occurred while loading comments due to " + loadError.Error()})
			return
		}

		context.JSON(http.StatusOK, comments)
	})


}
