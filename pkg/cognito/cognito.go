package cognito

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/gin-gonic/gin"
)

type AwsCognito struct {
	cognitoClient   *cognito.Client
	appClientId     string
	appClientSecret string
	appPoolId       string
}

func NewCognitoClient(config *aws.Config, clientId, clientSecret, poolId string) *AwsCognito {
	client := cognito.NewFromConfig(*config)

	return &AwsCognito{
		cognitoClient:   client,
		appClientId:     clientId,
		appClientSecret: clientSecret,
		appPoolId:       poolId,
	}
}

func computeSecretHash(clientSecret, username, clientId string) string {
	mac := hmac.New(sha256.New, []byte(clientSecret))
	mac.Write([]byte(username + clientId))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func (c *AwsCognito) SignUp(ctx *gin.Context, email, password string) (string, error) {
	input := &cognito.SignUpInput{
		ClientId: aws.String(c.appClientId),
		Username: aws.String(email),
		Password: aws.String(password),
	}

	secretHash := computeSecretHash(c.appClientSecret, email, c.appClientId)

	input.SecretHash = aws.String(secretHash)

	user, err := c.cognitoClient.SignUp(ctx, input)
	if err != nil {
		if strings.Contains(err.Error(), "UsernameExistsException") {
			err = errors.New("akun dengan email yang ingin didaftarkan telah ada, jika lupa password silahkan klik lupa password")
		}

		if strings.Contains(err.Error(), "InvalidParameterException") {
			err = errors.New("username harus berupa email")
		}
		return "", err
	}

	return *user.CodeDeliveryDetails.Destination, nil
}

func (c *AwsCognito) ConfirmSignUp(email, code string) (*cognito.ConfirmSignUpOutput, error) {
	input := &cognito.ConfirmSignUpInput{
		ClientId:         aws.String(c.appClientId),
		Username:         aws.String(email),
		ConfirmationCode: aws.String(code),
	}

	secretHash := computeSecretHash(c.appClientSecret, email, c.appClientId)
	input.SecretHash = aws.String(secretHash)

	result, err := c.cognitoClient.ConfirmSignUp(context.TODO(), input)
	if err != nil {
		if strings.Contains(err.Error(), "ExpiredCodeException") {
			err = errors.New("kode telah expire, silahkan request kode konfirmasi kembali")
		}

		if strings.Contains(err.Error(), "CodeMismatchException") {
			err = errors.New("kode verifikasi gagal, silahkan cek kembali pesan di email anda")
		}
		return nil, err
	}

	return result, nil
}

func (c *AwsCognito) SignIn(ctx *gin.Context, email, password string) (*types.AuthenticationResultType, error) {
	secretHash := computeSecretHash(c.appClientSecret, email, c.appClientId)

	input := &cognito.InitiateAuthInput{
		ClientId: aws.String(c.appClientId),
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		AuthParameters: map[string]string{
			"USERNAME":    email,
			"PASSWORD":    password,
			"SECRET_HASH": secretHash,
		},
	}

	result, err := c.cognitoClient.InitiateAuth(ctx, input)
	if err != nil {
		if strings.Contains(err.Error(), "NotAuthorizedException") {
			err = errors.New("name pengguna atau password kurang tepat")
		}

		if strings.Contains(err.Error(), "UserNotConfirmedException") {
			err = errors.New("akun belum terkonfirmasi, silahkan cek kode konfirmasi didalam email")
		}

		if strings.Contains(err.Error(), "InvalidParameterException") {
			err = errors.New("parameter tidak valid")
		}
		return nil, err
	}

	return result.AuthenticationResult, nil
}

func (c *AwsCognito) SignOut(ctx *gin.Context, username string) error {
	input := &cognito.AdminUserGlobalSignOutInput{
		UserPoolId: aws.String(c.appPoolId),
		Username:   aws.String(username),
	}
	_, err := c.cognitoClient.AdminUserGlobalSignOut(ctx, input)
	if err != nil {
		log.Println("failed to perform global sign out:", err.Error())
		return err
	}

	return nil
}

func (c *AwsCognito) ResendConfirmationCode(ctx *gin.Context, email string) (string, error) {
	input := &cognito.ResendConfirmationCodeInput{
		ClientId: aws.String(c.appClientId),
		Username: aws.String(email),
	}

	secretHash := computeSecretHash(c.appClientSecret, email, c.appClientId)
	input.SecretHash = aws.String(secretHash)

	code, err := c.cognitoClient.ResendConfirmationCode(ctx, input)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	return *code.CodeDeliveryDetails.Destination, nil
}

func (c *AwsCognito) GetUsername(ctx *gin.Context, token string) (*cognito.GetUserOutput, error) {
	input := &cognito.GetUserInput{
		AccessToken: &token,
	}
	result, err := c.cognitoClient.GetUser(ctx, input)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return result, nil
}

func (c *AwsCognito) UpdateUserAttributes(ctx *gin.Context, email, attribute, value string) error {
	inputAttributes := []types.AttributeType{
		{
			Name:  aws.String(attribute),
			Value: aws.String(value),
		},
	}

	input := &cognito.AdminUpdateUserAttributesInput{
		UserPoolId:     aws.String(c.appPoolId),
		Username:       aws.String(email),
		UserAttributes: inputAttributes,
	}

	_, err := c.cognitoClient.AdminUpdateUserAttributes(ctx, input)

	return err
}

func (c *AwsCognito) ForgotPassword(ctx *gin.Context, email string) (string, error) {
	input := &cognito.ForgotPasswordInput{
		ClientId: aws.String(c.appClientId),
		Username: aws.String(email),
	}

	secretHash := computeSecretHash(c.appClientSecret, email, c.appClientId)
	input.SecretHash = aws.String(secretHash)

	result, err := c.cognitoClient.ForgotPassword(ctx, input)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	return *result.CodeDeliveryDetails.Destination, nil
}

func (c *AwsCognito) ResetPassword(ctx *gin.Context, email, password, code string) (*cognito.ConfirmForgotPasswordOutput, error) {
	input := &cognito.ConfirmForgotPasswordInput{
		ClientId:         aws.String(c.appClientId),
		Username:         aws.String(email),
		Password:         aws.String(password),
		ConfirmationCode: aws.String(code),
	}

	secretHash := computeSecretHash(c.appClientSecret, email, c.appClientId)
	input.SecretHash = aws.String(secretHash)

	result, err := c.cognitoClient.ConfirmForgotPassword(ctx, input)
	if err != nil {
		if strings.Contains(err.Error(), "CodeMismatchException") {
			err = errors.New("kode verifikasi tidak cocok, silahkan request kode konfirmais kembali")
		}
		if strings.Contains(err.Error(), "LimitExceededException") {
			err = errors.New("maksimum pengulangan password hanya sampai 3x. Silahkan coba lagi untuk beberapa waktu")
		}
		if strings.Contains(err.Error(), "ExpiredCodeException") {
			err = errors.New("kode verifikasi gagal, silahkan request lupa password kembali")
		}
		if strings.Contains(err.Error(), "InvalidParameterException") {
			err = errors.New("parameter tidak bisa dipakai, silahkan ubah paramternya")
		}
		return nil, err
	}

	return result, err
}

func (c *AwsCognito) AddUserToGroup(ctx *gin.Context, username, groupName string) error {
	input := &cognito.AdminAddUserToGroupInput{
		UserPoolId: aws.String(c.appPoolId),
		Username:   aws.String(username),
		GroupName:  aws.String(groupName),
	}

	_, err := c.cognitoClient.AdminAddUserToGroup(ctx, input)
	if err != nil {
		if strings.Contains(err.Error(), "ResourceNotFoundException") {
			err = errors.New("group tidak ada")
		}

		//if strings.Contains(err.Error(), "UserNotFoundException") {
		//	err = errors.New("nama pengguna belum terdaftar")
		//}
		return err
	}
	return nil
}

func (c *AwsCognito) CheckUserInGroup(username string) ([]string, error) {
	input := &cognito.AdminListGroupsForUserInput{
		Username:   aws.String(username),
		UserPoolId: aws.String(c.appPoolId),
	}

	resp, err := c.cognitoClient.AdminListGroupsForUser(context.TODO(), input)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	groupNames := make([]string, len(resp.Groups))
	for i, group := range resp.Groups {
		groupNames[i] = *group.GroupName
	}

	return groupNames, nil
}

func (c *AwsCognito) ListGroup() ([]string, error) {
	input := &cognito.ListGroupsInput{
		UserPoolId: aws.String(c.appPoolId),
	}

	resp, err := c.cognitoClient.ListGroups(context.TODO(), input)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	groupNames := make([]string, len(resp.Groups))
	for i, group := range resp.Groups {
		groupNames[i] = *group.GroupName
	}
	return groupNames, nil
}
