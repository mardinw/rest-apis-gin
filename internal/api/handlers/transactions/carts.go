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

// @Summary CreateCart access process
// @Description do create a cart
// @Tags transactions
// @Accept json
// @Produce json
// @Param create body dtos.Carts true "create a cart"
// @Success 200 {object} dtos.MessagesResponses "the message successfully create"
// @Failure 400 {string} string "Error Bad Request"
// @Router /transactions/carts [post]
// @Security Bearer
func CreateCart(db *sql.DB) gin.HandlerFunc {
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

		if err := ctx.ShouldBindJSON(&cart); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

		if err := cart.Insert(*output.Username, db); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"message": "create successfully",
			"group":   resp,
		})
	}
}

// @Summary GetAll carts process
// @Description do get all carts
// @Tags transactions
// @Accept json
// @Produce json
// @Success 200 {object} dtos.MessagesResponses "get all transactions"
// @Failure 400 {string} string "not authorized"
// @Router /transactions/carts [get]
// @Security Bearer
func GetAllCart(db *sql.DB) gin.HandlerFunc {
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

		result, err := cart.GetAll(*output.Username, db)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"transactions": result})
	}
}

// @Summary GetID carts process
// @Description do get id carts
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path integer true "get id transaction"
// @Success 200 {object} dtos.MessagesResponses "get id transaction"
// @Failure 400 {string} string "unauthorized"
// @Router /transactions/carts/{id} [get]
// @Security Bearer
func GetIDCart(db *sql.DB) gin.HandlerFunc {
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
			ctx.JSON(http.StatusForbidden, gin.H{"message": "forbidden"})
			return
		}

		result, err := cart.GetID(int64(id), *output.Username, db)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"transaction": result})
	}
}

// @Summary update cart process
// @Description do update cart
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path integer true "access id transactions"
// @Param transaction body dtos.Carts true "update carts"
// @Success 200 {object} dtos.MessagesResponses "update transaction successfully"
// @Failure 400 {string} string "Error Bad request"
// @Router /transactions/carts/{id} [put]
// @Security Bearer
func UpdateChart(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var cart transactions.Carts

		authHeader := ctx.GetHeader("Authorization")
		splitted := strings.Split(authHeader, " ")

		token := splitted[1]
		output, err := aws.NewConnect().Cognito.GetUsername(ctx, token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		id, err := strconv.Atoi(ctx.Params.ByName("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		if err := ctx.ShouldBindJSON(&cart); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
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

		err = cart.Update(int64(id), *output.Username, db)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("update successfully id: %d", id),
		})
	}
}

// @Summary delete id transaction process
// @Description do delete a transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path integer true "delete charts"
// @Success 200 {object} dtos.MessagesResponses "delete successfully"
// @Failure 400 {string} string "unauthorized"
// @Router /transactions/carts/{id} [delete]
// @Security Bearer
func DeleteChart(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var cart transactions.Carts

		authHeader := ctx.GetHeader("Authorization")
		splitted := strings.Split(authHeader, " ")
		token := splitted[1]
		output, err := aws.NewConnect().Cognito.GetUsername(ctx, token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
		}
		id, err := strconv.Atoi(ctx.Params.ByName("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
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
		err = cart.Delete(int64(id), *output.Username, db)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("delete successfully id: %d", id),
		})
	}
}
