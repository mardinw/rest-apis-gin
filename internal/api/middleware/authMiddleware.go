package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sethvargo/go-envconfig"
	"payuoge.com/configs"
	"payuoge.com/dtos"
	"payuoge.com/internal/api/helpers"
	"payuoge.com/pkg/aws"
)

func Auth(redisClient *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var configs configs.AppConfiguration

		if err := envconfig.Process(ctx, &configs); err != nil {
			log.Fatal(err.Error())
		}

		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			ctx.Abort()
			return
		}

		splitted := strings.Split(authHeader, " ")
		if len(splitted) != 2 || strings.ToLower(splitted[0]) != "bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			ctx.Abort()
			return
		}

		accessToken := splitted[1]
		output, err := aws.NewConnect().Cognito.GetUsername(ctx, accessToken)
		if err != nil {
			log.Println(err.Error())
			return
		}
		key := fmt.Sprintf("user:%s:access_token", *output.Username)

		cachedToken, err := redisClient.Get(ctx, key).Result()
		if err != nil || cachedToken != accessToken {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// SignInHandler godoc
//
// @Summary Login access process
// @Description do login account
// @Tags auth
// @Accept json
// @Produce json
// @Param login body dtos.AuthData true "login data"
// @Success 200 {object} dtos.MessagesResponses "the message success responses"
// @Failure 400 {string} string "Error Bad request"
// @Router /auth/login [post]
func Login(redisClient *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var loginData dtos.AuthData
		if err := ctx.ShouldBindJSON(&loginData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result, err := aws.NewConnect().Cognito.SignIn(ctx, loginData.Username, loginData.Password)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		output, err := aws.NewConnect().Cognito.GetUsername(ctx, *result.AccessToken)
		if err != nil {
			log.Println(err.Error())
			return
		}

		key := fmt.Sprintf("user:%s:access_token", *output.Username)

		accessToken := *result.AccessToken

		err = redisClient.Set(ctx, key, accessToken, time.Hour).Err()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		cookie := http.Cookie{
			Name:     "access_token",
			Value:    accessToken,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			Expires:  time.Now().Add(time.Hour * 1),
		}
		http.SetCookie(ctx.Writer, &cookie)

		ctx.JSON(http.StatusOK, gin.H{
			"idToken": accessToken,
		})
	}
}

// @Summary Login google access process
// @Description do logout account
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} dtos.MessagesResponses "logged in via google successfully"
// @Failure 400 {string} string "cookie not found"
// @Router /auth/google [get]
func LoginGA(ctx *gin.Context) {
	var configs configs.AppConfiguration
	if err := envconfig.Process(context.Background(), &configs); err != nil {
		log.Fatal(err.Error())
	}

	URI := "https://auth.payuoge.com/oauth2/authorize?"
	provider := "Google"
	response := "CODE"
	clientID := configs.AwsConf.ClientId
	scope := "aws.cognito.signin.user.admin+email+openid"

	if configs.AppEnv == "production" {
		redirectURI := "https://api.payuoge.com/v1/auth/callback"
		accessURI := URI +
			"identity_provider=" + provider +
			"&redirect_uri=" + redirectURI +
			"&response_type=" + response +
			"&client_id=" + clientID +
			"&scope=" + scope
		ctx.Redirect(http.StatusMovedPermanently, accessURI)
	} else {
		redirectURI := "http://localhost:4001/v1/auth/callback"
		accessURI := URI +
			"identity_provider=" + provider +
			"&redirect_uri=" + redirectURI +
			"&response_type=" + response +
			"&client_id=" + clientID +
			"&scope=" + scope
		ctx.Redirect(http.StatusMovedPermanently, accessURI)
	}
}

// LogoutHandler godoc
//
// @Summary Logout access process
// @Description do logout account
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} dtos.MessagesResponses "logged out successfully"
// @Failure 400 {string} string "cookie not found"
// @Router /auth/logout [get]
// @Security Bearer
func Logout(redisClient *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Request.Cookie("access_token")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Cookie not found"})
			return
		}
		output, err := aws.NewConnect().Cognito.GetUsername(ctx, cookie.Value)
		if err != nil {
			log.Println(err.Error())
			return
		}
		key := fmt.Sprintf("user:%s:access_token", *output.Username)

		err = redisClient.Del(ctx, key).Err()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed deleted token from redis"})
			return
		}

		expiredCookie := http.Cookie{
			Name:     "access_token",
			Value:    "",
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			Expires:  time.Now().Add(-1 * time.Hour),
		}
		http.SetCookie(ctx.Writer, &expiredCookie)
		ctx.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
	}
}

// RegisterHandler godoc
//
// @Summary Register access process
// @Description do register account
// @Tags auth
// @Accept json
// @Produce json
// @Param registerData body dtos.AuthData true "used for register"
// @Success 200 {object} dtos.MessagesResponses "send code confirmation"
// @Failure 400 {string} string "Error Bad request"
// @Router /auth/register [post]
func Register(ctx *gin.Context) {
	var registerData dtos.AuthData

	if err := ctx.ShouldBindJSON(&registerData); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	result, err := aws.NewConnect().Cognito.SignUp(ctx, registerData.Username, registerData.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Kode konfirmasi telah dikirim ke %s.", result),
	})
}

func CallbackCognito(redisClient *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var configs configs.AppConfiguration
		if err := envconfig.Process(context.Background(), &configs); err != nil {
			log.Fatal(err.Error())
		}

		code := ctx.DefaultQuery("code", "")
		if code == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing code parameter"})
			return
		}

		// configuration oauth2
		clientID := configs.AwsConf.ClientId
		clientSecret := configs.AwsConf.ClientSecret
		grantType := "authorization_code"
		redirectURI := configs.GoogleAuth.RedirectURL

		data := url.Values{}
		data.Set("grant_type", grantType)
		data.Set("client_id", clientID)
		data.Set("client_secret", clientSecret)
		data.Set("code", code)
		data.Set("redirect_uri", redirectURI)

		resp, err := http.Post("https://auth.payuoge.com/oauth2/token", "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))

		if err != nil {
			log.Println(err.Error())
		}

		defer resp.Body.Close()

		var idToken map[string]interface{}

		if err := json.NewDecoder(resp.Body).Decode(&idToken); err != nil {
			return
		}

		output, err := aws.NewConnect().Cognito.GetUsername(ctx, idToken["access_token"].(string))
		if err != nil {
			log.Println(err.Error())
			return
		}

		key := fmt.Sprintf("user:%s:access_token", *output.Username)
		err = redisClient.Set(ctx, key, idToken["access_token"].(string), time.Hour).Err()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		cookie := http.Cookie{
			Name:     "access_token",
			Value:    idToken["access_token"].(string),
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			Expires:  time.Now().Add(time.Hour * 1),
		}
		http.SetCookie(ctx.Writer, &cookie)

		ctx.Redirect(http.StatusMovedPermanently, "https://payuoge.com")
	}
}

// ConfirmationHandler godoc
//
// @Summary Confirmation SignUp access process
// @Description do confirmation account
// @Tags auth
// @Accept json
// @Produce json
// @Param confirmData body dtos.AuthCodeData true "used for confirmation an account"
// @Success 200 {object} string "redirect to home"
// @Failure 400 {string} string "Error Bad request"
// @Router /auth/confirm [post]
func Confirmation(ctx *gin.Context) {
	var confirmData dtos.AuthCodeData

	if err := ctx.BindJSON(&confirmData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := aws.NewConnect().Cognito.ConfirmSignUp(confirmData.Username, confirmData.ConfirmationCode)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "akun telah dikonfirmasi"})
}

// ResendCodeHandler godoc
//
// @Summary ResendCode process
// @Description do resend code to an account
// @Tags auth
// @Accept json
// @Produce json
// @Param resendCode body dtos.Users true "used for resend code"
// @Success 200 {object} dtos.MessagesResponses "check inbox email address"
// @Failure 400 {string} string "Error Bad request"
// @Router /auth/resend [post]
func ResendCode(ctx *gin.Context) {
	var resendCode dtos.Users

	if err := ctx.BindJSON(&resendCode); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := aws.NewConnect().Cognito.ResendConfirmationCode(ctx, resendCode.Username)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("kode konfirmasi telah dikirim ulang kembali ke %s", result)})
}

// @Summary ForgotPassword process
// @Description do forgot password to an account
// @Tags auth
// @Accept json
// @Produce json
// @Param forgotData body dtos.Users true "used for forgot password"
// @Success 200 {object} dtos.MessagesResponses "send code confirmation for forgot password"
// @Failure 400 {string} string "Error Bad request"
// @Router /auth/forgot [post]
func ForgotPassword(ctx *gin.Context) {
	var forgotData dtos.Users

	if err := ctx.BindJSON(&forgotData); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	result, err := aws.NewConnect().Cognito.ForgotPassword(ctx, forgotData.Username)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("kode konfirmasi untuk lupa password telah dikirim, silahkan cek inbox %s", result)})
}

// ResetPasswordHandler godoc
//
// @Summary ResetPassword process
// @Description do reset password to an account
// @Tags auth
// @Accept json
// @Produce json
// @Param resetPassword body dtos.AuthResetData true "used for reset password"
// @Success 200 {object} dtos.MessagesResponses "send for reset password"
// @Failure 400 {string} string "Error Bad request"
// @Router /auth/reset [post]
func ResetPassword(ctx *gin.Context) {
	var resetPassword dtos.AuthResetData
	if err := ctx.BindJSON(&resetPassword); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	_, err := aws.NewConnect().Cognito.ResetPassword(ctx, resetPassword.Username, resetPassword.Password, resetPassword.ConfirmationCode)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "reset password telah sukses"})
}

// Add User To Group godoc
//
// @Summary Add user to group process
// @Description do add user to group
// @Tags auth
// @Accept json
// @Produce json
// @Param addUserToGroup body dtos.UpdateGroup true "user add to group"
// @Success 200 {object} dtos.MessagesResponses "user success to add group"
// @Failure 400 {string} string "Error Bad request"
// @Router /auth/add-user-groups [post]
// @Security Bearer
func AddUserToGroup(redisClient *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var addUserToGroup dtos.UpdateGroup
		authHeader := ctx.GetHeader("Authorization")
		splitted := strings.Split(authHeader, " ")

		token := splitted[1]

		output, err := aws.NewConnect().Cognito.GetUsername(ctx, token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		if err := ctx.BindJSON(&addUserToGroup); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		err = aws.NewConnect().Cognito.AddUserToGroup(ctx, *output.Username, addUserToGroup.Groups)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("pengguna telah ditambahkan ke group %s", addUserToGroup.Groups)})
	}

}

// @Summary list group process
// @Description do list groups
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} dtos.MessagesResponses "list group"
// @Failure 400 {string} string "Error Bad request"
// @Router /auth/list-groups [get]
// @Security Bearer
func GetGroups() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		splitted := strings.Split(authHeader, " ")

		token := splitted[1]

		output, err := aws.NewConnect().Cognito.GetUsername(ctx, token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		err = helpers.CheckAccountAdmin(output.Username)
		if err != nil {
			ctx.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
			return
		}

		resp, err := aws.NewConnect().Cognito.ListGroup()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"groups": resp})
	}
}
