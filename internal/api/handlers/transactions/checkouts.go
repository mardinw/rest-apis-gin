package transactions

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"payuoge.com/internal/api/models/transactions"
	"payuoge.com/pkg/aws"
)

// @Summary CreateCheckout access process
// @Description do create a create checkout
// @Tags transactions
// @Accept json
// @Produce json
// @Success 200 {object} dtos.MessagesResponses "message responses"
// @Failure 400 {string} string "not authorized"
// @Router /transactions/checkout [post]
// @Security Bearer
func CreateCheckout(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var checkout transactions.Checkouts
		// var cart transactions.Carts

		authHeader := ctx.GetHeader("Authorization")
		splitted := strings.Split(authHeader, " ")
		token := splitted[1]
		output, err := aws.NewConnect().Cognito.GetUsername(ctx, token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
		}

		// check roles groups
		resp, err := aws.NewConnect().Cognito.CheckUserInGroup(*output.Username)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		targetValue := "retail"
		found := false
		for _, value := range resp {
			if value == targetValue {
				found = true
				break
			}
		}

		if found {
			log.Printf("Found %s\n", targetValue)
		} else {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "permission denied"})
			return
		}

		totalAmount := checkout.CalculateTotalAmount(*output.Username, db)
		err = checkout.Insert(*output.Username, totalAmount, db)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "checkout successfully"})
	}
}

// @Summary Get All Checkout access process
// @Description do Get All checkout
// @Tags transactions
// @Accept json
// @Produce json
// @Success 200 {object} dtos.MessagesResponses "message responses"
// @Failure 400 {string} string "not authorized"
// @Router /transactions/checkout [get]
// @Security Bearer
func GetCheckout(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var checkout transactions.Checkouts

		authHeader := ctx.GetHeader("Authorization")
		splitted := strings.Split(authHeader, " ")
		token := splitted[1]
		output, err := aws.NewConnect().Cognito.GetUsername(ctx, token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
		}

		// check roles groups
		resp, err := aws.NewConnect().Cognito.CheckUserInGroup(*output.Username)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		targetValue := "retail"
		found := false
		for _, value := range resp {
			if value == targetValue {
				found = true
				break
			}
		}

		if found {
			log.Printf("Found %s\n", targetValue)
		} else {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "permission denied"})
			return
		}

		result, err := checkout.GetAll(*output.Username, db)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"checkouts": result})
	}
}

// @Summary Get id Checkout access process
// @Description do Get id checkout
// @Tags transactions
// @Accept json
// @Produce json
// @Success 200 {object} dtos.MessagesResponses "message responses"
// @Failure 400 {string} string "not authorized"
// @Router /transactions/checkout/{id} [get]
// @Security Bearer
func GetIDCheckout(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var cart transactions.Carts

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

		id, err := strconv.Atoi(ctx.Params.ByName("id"))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// check roles groups
		resp, err := aws.NewConnect().Cognito.CheckUserInGroup(*output.Username)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		targetValue := "retail"
		found := false
		for _, value := range resp {
			if value == targetValue {
				found = true
				break
			}
		}

		if found {
			log.Printf("Found %s\n", targetValue)
		} else {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "permission denied"})
			return
		}

		result, err := cart.GetID(int64(id), *output.Username, db)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"checkouts": result})
	}
}

// @Summary update id Checkout access process
// @Description do update id checkout
// @Tags transactions
// @Accept json
// @Produce json
// @Success 200 {object} dtos.MessagesResponses "message responses"
// @Failure 400 {string} string "not authorized"
// @Router /transactions/checkout/{id} [put]
// @Security Bearer
func UpdateCheckout(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var checkout transactions.Checkouts

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

		id, err := strconv.Atoi(ctx.Params.ByName("id"))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// check roles groups
		resp, err := aws.NewConnect().Cognito.CheckUserInGroup(*output.Username)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		targetValue := "retail"
		found := false
		for _, value := range resp {
			if value == targetValue {
				found = true
				break
			}
		}

		if found {
			log.Printf("Found %s\n", targetValue)
		} else {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "permission denied"})
			return
		}

		err = checkout.Update(int64(id), *output.Username, db)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("checkout update successfully id: %d", id),
		})
	}
}

// @Summary delete id Checkout access process
// @Description do delete id checkout
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path integer true "delete by id"
// @Success 200 {object} dtos.MessagesResponses "message responses"
// @Failure 400 {string} string "not authorized"
// @Router /transactions/checkout/{id} [delete]
// @Security Bearer
func DeleteCheckout(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var checkout transactions.Checkouts

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

		id, err := strconv.Atoi(ctx.Params.ByName("id"))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// check roles groups
		resp, err := aws.NewConnect().Cognito.CheckUserInGroup(*output.Username)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		targetValue := "retail"
		found := false
		for _, value := range resp {
			if value == targetValue {
				found = true
				break
			}
		}

		if found {
			log.Printf("Found %s\n", targetValue)
		} else {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "forbidden"})
			return
		}

		err = checkout.Delete(int64(id), *output.Username, db)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("delete successfully checkout :%d", id),
		})
	}
}
