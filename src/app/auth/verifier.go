package auth

import (
	"errors"
	"fmt"
	"go-social-app/src/app/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"

)

type OauthVerifier struct {

}

func (*OauthVerifier) ValidateUser(username, password, scope string, req *http.Request) error  {

	user := models.User{}
	result := models.Db.Where(&models.User{Username: username}).First(&user)

	if result.Error != nil {
		return result.Error
	}

	if err:= bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(password)); err != nil {
		return errors.New("passwords do not match")
	}

	return nil

}

func (*OauthVerifier) ValidateClient(clientID, clientSecret, scope string, req *http.Request) error  {

	oauthClientDetails := models.OauthClientDetails{}
	result := models.Db.Where(&models.OauthClientDetails{ClientId: clientID}).First(&oauthClientDetails)

	if result.Error != nil {
		return errors.New("client not found")
	}

	err := result.Scan(&oauthClientDetails.ClientSecret)

	if err!= nil {
		return errors.New("invalid client secret")
	}

	if err:= bcrypt.CompareHashAndPassword([]byte(oauthClientDetails.ClientSecret),[]byte(clientSecret)); err != nil {
		return errors.New("invalid client secret")
	}

	return nil
}

func (*OauthVerifier) AddClaims(credential, tokenID, tokenType, scope string) (map[string]string, error) {
	claims:=map[string]string{}
	return claims, nil
}

func (*OauthVerifier) StoreTokenId(credential, tokenID, refreshTokenID, tokenType string) error {
	return nil
}

func (*OauthVerifier) AddProperties(credential, tokenID, tokenType string, scope string) (map[string]string, error)  {
	fmt.Println("Credential is ", credential)
	user := models.User{}
	models.Db.Where(&models.User{Username: credential}).First(&user)
	properties := map[string]string{
		"username" : user.Username,
		"first_name" : user.FirstName,
		"last_name" : user.LastName,
		"gender" : user.Gender,
		"display_picture" : user.DisplayPicture,
		"email" : user.Email,
	}

	return properties, nil
}

func (*OauthVerifier) ValidateTokenId(credential, tokenID, refreshTokenID, tokenType string) error  {
	return nil
}