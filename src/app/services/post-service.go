package services

import "go-social-app/src/app/models"

func LoadPostsWithLesserId(id int, pageSize int) ([]models.Post, error)  {

	var posts []models.Post
	result := models.Db.Order("id desc").Preload("User").Where("id < ? ", id).Limit(pageSize).Find(&posts)

	if result.Error != nil {
		return nil, result.Error
	}

	return posts, nil

}

func LoadPostByPageNumberAndSize(pageNo int, pageSize int) ([]models.Post, error) {
	var posts []models.Post
	result := models.Db.Order("id desc").Preload("User").Offset(pageNo*pageSize).Limit(pageSize).Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}

	return posts, nil

}