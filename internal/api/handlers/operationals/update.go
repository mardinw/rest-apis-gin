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

// @Summary UpdateOperational access process
// @Description do Update a operational
// @Tags groceries
// @Accept json
// @Product json
// @Param id path integer true "id operational"
// @Param operational body dtos.Operate true "update a operational"
// @Success 200 {object} dtos.MessagesResponses "the message successfully create"
// @Failure 400 {string} string "Error Bad Request"
// @Router /groceries/operational/{id} [put]
// @Security Bearer
func Update(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var operateData models.Operationals
		var updateData models.Operationals

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
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if err := ctx.ShouldBindJSON(&operateData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		if updateData.DayOperational != nil {
			operateData.DayOperational = updateData.DayOperational
		}

		if updateData.Open != 0 {
			operateData.Open = updateData.Open
		}

		if updateData.Close != 0 {
			operateData.Close = updateData.Close
		}

		if updateData.Active != false {
			updateData.Active = true
		}

		// check roles group
		resp, err := aws.NewConnect().Cognito.CheckUserInGroup(*output.Username)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		targetValue := []string{
			"GROSIR_MANAGER", "GROSIR_OWNER",
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

		err = operateData.Update(uint64(id), *output.Username, db)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"messsage": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("update id %d successfully", id),
		})
	}
}
