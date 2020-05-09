package services

import (
	"go-social-app/src/app/models"
)

func LoadUserByUsername(username string) (models.User, error) {
	var user = models.User{}
	userResult := models.Db.Where(&models.User{Username: username}).First(&user)

	return user, userResult.Error

}

func FilterWithSearchValue(filter string) ([]models.User, error) {

	var users = []models.User{}
	userResult := models.Db.Where("Username like ?",
		"%"+filter+"%").Or("First_Name like ?",
			"%"+filter+"%").Or("Last_Name like ?", "%"+filter+"%").Find(&users)

	return users, userResult.Error
}
