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

// @Summary DeleteID product process
// @Description do delete a product
// @Tags groceries
// @Accept json
// @Produce json
// @Param id path integer true "delete a product"
// @Success 200 {object} dtos.MessagesResponses "delete a product"
// @Failure 400 {string} string "cookie not found"
// @Router /groceries/products/{id} [delete]
// @Security Bearer
func DeleteID(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var product products.Product

		// check authorized
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

		// check roles group
		err = helpers.CheckAccountGroceries(output.Username)
		if err != nil {
			ctx.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
			return
		}

		err = product.Delete(int64(id), *output.Username, db)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("delete successfully id:%d", id),
		})
	}
}
