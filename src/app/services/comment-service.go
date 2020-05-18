package services

import "go-social-app/src/app/models"

func LoadCommentsByPostIdAndPageNoAndSize(postId, pageNo, pageSize int) ([]models.Comment, error) {
	comments := []models.Comment{}
	result := models.Db.Order("id desc").Preload("User").Where("postid = ? ", postId).Offset(pageNo*pageSize).Limit(pageSize).Find(&comments)

	return comments, result.Error
}

func LoadCommentsWithLesserId(id, postId, pageSize int) ([]models.Comment, error)  {

	var comments []models.Comment
	result := models.Db.Order("id desc").Preload("User").Where("id < ? and postid = ? ", id, postId).Limit(pageSize).Find(&comments)

	if result.Error != nil {
		return nil, result.Error
	}

	return comments, nil

}
