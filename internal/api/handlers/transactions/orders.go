package transactions

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"payuoge.com/internal/api/helpers"
	"payuoge.com/internal/api/models/transactions"
	"payuoge.com/pkg/aws"
)

// @Summary Create Orders access process
// @Description do create a order
// @Tags transactions
// @Accept json
// @Produce json
// @Router /transactions/orders [post]
// @Security Bearer
func CreateOrders(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var order transactions.Orders
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
		err = helpers.CheckAccountRetail(output.Username)
		if err != nil {
			ctx.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
			return
		}

		totalAmount := checkout.CalculateTotalAmount(*output.Username, db)

		err = order.Insert(*output.Username, int32(totalAmount), db)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "order success"})
	}
}

// @Summary Get All Order access process
// @Description do get all order
// @Tags transactions
// @Accept json
// @Produce json
// @Router /transactions/orders [get]
// @Security Bearer
func GetOrders(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var order transactions.Orders

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

		result, err := order.GetAll(*output.Username, db)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"orders": result})
	}
}

// @Summary Get id order access process
// @Description do get id order
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path integer true "get id order"
// @Router /transactions/orders/{id} [get]
// @Security Bearer
func GetIDOrder(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var order transactions.Orders

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

		result, err := order.GetID(int64(id), *output.Username, db)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "not found record",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"order": result})
	}
}

// @Summary update id order access process
// @Description do update id order
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path integer true "update id order"
// @Router /transactions/orders/{id} [put]
// @Security Bearer
func UpdateOrder(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var order transactions.Orders

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

		err = order.Update(int64(id), *output.Username, db)
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

// @Summary delete id order access process
// @Description do delete id order
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path integer true "delete by id"
// @Router /transactions/orders/{id} [delete]
// @Security Bearer
func DeleteOrder(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var order transactions.Orders

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

		err = order.Delete(int64(id), *output.Username, db)
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
