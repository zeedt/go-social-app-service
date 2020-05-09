package s3

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"mime/multipart"
)

func UploadFilesToS3(files []*multipart.FileHeader) ([]string, error)  {

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
		Credentials: credentials.NewEnvCredentials(),
	}))
	uploader := s3manager.NewUploader(sess)

	var filePaths []string
	for _, file := range files {
		f, err := file.Open()
		if err != nil {
			fmt.Println(err.Error())
			return nil, errors.New(err.Error())
		}
		uploadResult, err := uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String("social-app-bucket1"),
			Key:    aws.String("social-app-images/" + file.Filename),
			Body:   f,
		})

		filePaths = append(filePaths, uploadResult.Location)

	}

	return filePaths, nil
}
