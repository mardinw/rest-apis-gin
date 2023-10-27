package operationals

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"payuoge.com/internal/api/models"
	"payuoge.com/pkg/aws"
)

// @Summary CreateOperational access process
// @Description do create a operational
// @Tags groceries
// @Accept json
// @Produce json
// @Param operational body dtos.Operate true "create a operational"
// @Success 200 {object} dtos.MessagesResponses "the message successfully create"
// @Failure 400 {string} string "Error Bad Request"
// @Router /groceries/operational [post]
// @Security Bearer
func Create(db *sql.DB) gin.HandlerFunc {
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
			return
		}

		if err := ctx.ShouldBindJSON(&operate); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

		if err := operate.Insert(*output.Username, db); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"message": "create operationals successfully",
		})
	}
}
