package cognito

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/sethvargo/go-envconfig"
	"golang.org/x/oauth2"
	"payuoge.com/configs"
)

type GoogleAuthenticator struct {
	CognitoClient    *cognito.Client
	GoogleAuthConfig *oauth2.Config
}

func NewGoogleAuthenticator(cfg *aws.Config) *GoogleAuthenticator {

	var configs configs.AppConfiguration

	if err := envconfig.Process(context.Background(), &configs); err != nil {
		log.Fatal(err.Error())
	}

	cognitoClient := cognito.NewFromConfig(*cfg)

	googleOAuthConfig := &oauth2.Config{
		ClientID:     configs.GoogleAuth.ClientID,
		ClientSecret: configs.GoogleAuth.ClientSecret,
		RedirectURL:  configs.GoogleAuth.RedirectURL,
		Scopes:       []string{"openid", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
	}

	return &GoogleAuthenticator{
		CognitoClient:    cognitoClient,
		GoogleAuthConfig: googleOAuthConfig,
	}
}

// func (ga *GoogleAuthenticator) CognitoExchange(accessToken, email string) (*cognito.InitiateAuthOutput, error) {
// 	var configs configs.AppConfiguration
//
// 	if err := envconfig.Process(context.Background(), &configs); err != nil {
// 		log.Fatal(err.Error())
// 	}
//
// 	secretHash := computeSecretHash(configs.AwsConf.ClientSecret, email, configs.AwsConf.ClientId)
// 	input := &cognito.credential{
// 		ClientId: &configs.AwsConf.ClientId,
// 		AuthFlow: types.AuthFlowTypeUserSrpAuth,
// 		AuthParameters: map[string]string{
// 			"USERNAME":    email,
// 			"SRP_A":       accessToken,
// 			"SECRET_HASH": secretHash,
// 		},
// 	}
//
// 	result, err := ga.CognitoClient.InitiateAuth(context.Background(), input)
// 	if err != nil {
// 		log.Println(err)
// 	}
//
// 	return result, nil
// }

func (ga *GoogleAuthenticator) GetUserInfoFromGoogle(accessToken string) (map[string]interface{}, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}

func (ga *GoogleAuthenticator) GenerateRandomKey(length int) (string, error) {
	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(key), nil
}
