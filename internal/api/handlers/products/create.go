package products

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"payuoge.com/internal/api/helpers"
	"payuoge.com/internal/api/models/products"
	"payuoge.com/pkg/aws"
)

// CreateProductHandler godoc
//
// @Summary CreateProduct access process
// @Description do create a product
// @Tags groceries
// @Accept json
// @Produce json
// @Param product body dtos.Product true "create a product"
// @Success 200 {object} dtos.MessagesResponses "the message successfully create"
// @Failure 400 {string} string "Error Bad Request"
// @Router /groceries/products [post]
// @Security Bearer
func Create(db *sql.DB) gin.HandlerFunc {
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

		if err := ctx.ShouldBindJSON(&product); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// check roles group
		err = helpers.CheckAccountGroceries(output.Username)
		if err != nil {
			ctx.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
			return
		}

		if err := product.Insert(db, product.SizeTypeId, product.CategoryId, *output.Username); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"message": "create successfully",
		})
	}
}
