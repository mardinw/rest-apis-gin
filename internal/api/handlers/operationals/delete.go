package operationals

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"payuoge.com/internal/api/models"
	"payuoge.com/pkg/aws"
)

// @Summary delete operational process
// @Description do delete operational
// @Tags groceries
// @Accept json
// @Produce json
// @Param id path integer true "get id operational"
// @Success 200 {object} dtos.MessagesResponses "delete operational successfully"
// @Failure 400 {string} string "Error Bad request"
// @Router /groceries/operational/{id} [delete]
// @Security Bearer
func DeleteID(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var operate models.Operationals

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
		resp, err := aws.NewConnect().Cognito.CheckUserInGroup(*output.Username)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		targetValue := []string{
			"grosir",
		}

		found := false
		for _, value := range resp {
			for _, tV := range targetValue {
				if value == tV {
					found = true
					break
				}
			}
		}

		if found {
			log.Printf("Found %s\n", targetValue)
		} else {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "forbidden"})
			return
		}

		err = operate.Delete(uint64(id), *output.Username, db)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		ctx.JSON(http.StatusAccepted, gin.H{
			"message": fmt.Sprintf("delete id %d successfully", id),
		})
	}
}
