package products

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"payuoge.com/internal/api/models/products"
	"payuoge.com/pkg/aws"
)

// @Summary GetAll product process
// @Description do get all product
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {object} dtos.MessagesResponses "get all category"
// @Failure 400 {string} string "cookies not found"
// @Router /products [get]
// @Security Bearer
func GetAll(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var product products.Product

		authHeader := ctx.GetHeader("Authorization")
		splitted := strings.Split(authHeader, " ")

		token := splitted[1]

		_, err := aws.NewConnect().Cognito.GetUsername(ctx, token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		result, err := product.GetAll(db)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"products": result})
	}
}

// @Summary GetAll product groceries process
// @Description do get all product
// @Tags groceries
// @Accept json
// @Produce json
// @Router /groceries/products [get]
// @Security Bearer
func GetProductsGrocery(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var product products.Product

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

		result, err := product.GetProductsGroceries(*output.Username, db)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"products": result})
	}
}

// @Summary GetID product process
// @Description do get id product
// @Tags products
// @Accept json
// @Produce json
// @Param id path integer true "get id product"
// @Success 200 {object} dtos.MessagesResponses "get id category"
// @Failure 400 {string} string "cookies not found"
// @Router /product/{id} [get]
// @Security Bearer
func GetID(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var product products.Product

		id, err := strconv.Atoi(ctx.Params.ByName("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		authHeader := ctx.GetHeader("Authorization")
		splitted := strings.Split(authHeader, " ")

		token := splitted[1]

		_, err = aws.NewConnect().Cognito.GetUsername(ctx, token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		result, err := product.Get(int64(id), db)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"product": result})
	}
}
