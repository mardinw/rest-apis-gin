package products

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"payuoge.com/internal/api/helpers"
	"payuoge.com/internal/api/models/products"
	"payuoge.com/pkg/aws"
)

// @Summary UpdateProduct access process
// @Description do update a product
// @Tags groceries
// @Accept json
// @Produce json
// @Param id path integer true "id a product"
// @Param product body dtos.Product true "update a product"
// @Success 200 {object} dtos.MessagesResponses "the message successfully create"
// @Failure 400 {string} string "Error Bad Request"
// @Failure 404 {string} string "Not Found"
// @Router /groceries/products/{id} [put]
// @Security Bearer
func Update(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var product products.Product
		var updateData products.Product

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
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		productData, err := product.Get(int64(id), db)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}

		if err := ctx.ShouldBindJSON(&updateData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		if updateData.ProductCode != "" {
			productData.ProductCode = updateData.ProductCode
		}

		if updateData.ProductName != "" {
			productData.ProductName = updateData.ProductName
		}

		if updateData.Picture != "" {
			productData.Picture = updateData.Picture
		}

		if updateData.Position != "" {
			productData.Position = updateData.Position
		}

		if updateData.SizeTypeId != 0 {
			productData.SizeTypeId = updateData.SizeTypeId
		}

		if updateData.CategoryId != 0 {
			productData.CategoryId = updateData.CategoryId
		}

		if updateData.MRP != int32(0) {
			productData.MRP = updateData.MRP
		}

		if updateData.BuyPrice != int32(0) {
			productData.BuyPrice = updateData.BuyPrice
		}

		if updateData.Defective != int32(0) {
			productData.Defective = updateData.Defective
		}

		if !updateData.Active {
			productData.Active = true
		}

		// check roles group
		err = helpers.CheckAccountGroceries(output.Username)
		if err != nil {
			ctx.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
			return
		}

		err = productData.Update(db, productData.ID, *output.Username)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("product %s successfully updated", productData.ProductName),
		})
	}
}
