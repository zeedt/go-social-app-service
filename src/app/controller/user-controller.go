package controller

import (
	"github.com/gin-gonic/gin"
	"go-social-app/src/app/models"
	"go-social-app/src/app/s3"
	"go-social-app/src/app/services"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func InitiateUserController(route *gin.Engine, securedRoute *gin.RouterGroup) {

	route.POST("/signup", func(context *gin.Context) {
		var user models.User
		err := context.BindJSON(&user)

		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"message" : "Bad request",
			})
			return
		}
		passwordBytes := []byte(user.Password)
		hashedPassword, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
		if err != nil {
			panic(err.Error())
		}
		user.Password = string(hashedPassword)
		result := models.Db.Create(&user)
		if result.Error != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"message" : result.Error.Error(),
			})
			return
		} else {
			context.JSON(http.StatusOK, gin.H{
				"message" : "Successfully profiled user",
			})
		}

	})

	securedRoute.POST("/hello", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"message" : "Hello"})
	})

	securedRoute.GET("/my-info", func(context *gin.Context) {
		username , _ := context.Get("oauth.credential")
		user, error := services.LoadUserByUsername(username.(string))
		if error != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message" : "No user found"})
			return
		}

		context.JSON(http.StatusOK, user)
	})

	securedRoute.POST("/upload-image", func(context *gin.Context) {
		form, err := context.MultipartForm()
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message" : err.Error()})
		}

		files :=form.File["file"]
		var filePaths []string
		if len(files) > 0 {
			filePaths, err = s3.UploadFilesToS3(files)
		}
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message" : err.Error()})
			return
		}
		username , _ := context.Get("oauth.credential")
		var user=models.User{}
		models.Db.Model(&user).Where("username = ? ", username.(string)).Update("display_picture", filePaths[0])
		context.JSON(http.StatusOK, gin.H{"display_picture":filePaths[0]})
	})

	securedRoute.GET("/filter-users/:value", func(context *gin.Context) {
		searchValue := context.Param("value")

		users, err := services.FilterWithSearchValue(searchValue)

		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message" : err.Error()})
			return
		}

		context.JSON(http.StatusOK, users)

	})

	securedRoute.GET("/find-user/:value", func(context *gin.Context) {
		searchValue := context.Param("value")

		user, err := services.LoadUserByUsername(searchValue)

		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message" : err.Error()})
			return
		}

		context.JSON(http.StatusOK, user)

	})

	securedRoute.POST("/update-password", func(context *gin.Context) {
		var profileUpdate models.ProfileUpdate
		context.BindJSON(&profileUpdate)

		username , _ := context.Get("oauth.credential")
		user, error := services.LoadUserByUsername(username.(string))
		if error != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message" : error})
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(profileUpdate.CurrentPassword)); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message" : "Invalid current password"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(profileUpdate.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message" : "New password invalid"})
			return
		}
		models.Db.Model(&user).Where("username = ? ", user.Username).Update("password", string(hashedPassword))
		context.JSON(http.StatusOK,gin.H{"success" : true, "message" : "User's password successfully updated"})
	})

	securedRoute.POST("/update-my-info", func(context *gin.Context) {
		var profileUpdate models.ProfileUpdate
		context.BindJSON(&profileUpdate)

		username , _ := context.Get("oauth.credential")
		user, error := services.LoadUserByUsername(username.(string))
		if error != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message" : error})
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(profileUpdate.Password)); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message" : "Invalid current password"})
			return
		}
		if len(profileUpdate.Email) < 10 {
			context.JSON(http.StatusBadRequest, gin.H{"message" : "Email cannot be less than 10 characters"})
			return
		}
		if len(profileUpdate.FirstName) < 3 {
			context.JSON(http.StatusBadRequest, gin.H{"message" : "First name cannot be less than 3 characters"})
			return
		}
		if len(profileUpdate.LastName) < 3 {
			context.JSON(http.StatusBadRequest, gin.H{"message" : "Last name cannot be less than 3 characters"})
			return
		}

		models.Db.Model(&user).Where("username = ? ", user.Username).Update(
			"first_name", profileUpdate.FirstName).Update("last_name",
				profileUpdate.LastName).Update("email", profileUpdate.Email)

		context.JSON(http.StatusOK, gin.H{"successful": true, "message": "User information updated successfully"})

	})

}
