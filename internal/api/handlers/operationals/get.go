package operationals

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"payuoge.com/internal/api/models"
	"payuoge.com/pkg/aws"
)

// @Summary GetAllOperational access process
// @Description do Get all operational
// @Tags groceries
// @Accept json
// @Produce json
// @Success 200 {object} dtos.MessagesResponses "the message successfully create"
// @Failure 400 {string} string "Error Bad Request"
// @Router /groceries/operational [get]
// @Security Bearer
func GetAll(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var operate models.Operationals

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

		result, err := operate.GetAll(db)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"operationals": result})
	}
}

// @Summary GetIDOperational access process
// @Description do Get ID operational
// @Tags groceries
// @Accept json
// @Product json
// @Param id path integer true "id operational"
// @Router /groceries/operational/{id} [get]
// @Security Bearer
func GetId(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var operate models.Operationals
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

		id, err := strconv.Atoi(ctx.Params.ByName("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result, err := operate.Get(uint64(id), db)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"operational": result,
		})
	}
}
