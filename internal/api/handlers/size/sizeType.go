package size

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"payuoge.com/internal/api/models/products"
	"payuoge.com/pkg/aws"
)

// @Summary create sizeType process
// @Description do create size type
// @Tags products
// @Accept json
// @Produce json
// @Param size body dtos.SizeType true "create sizetype"
// @Success 200 {object} dtos.MessagesResponses "create size type successfully"
// @Failure 400 {string} string "Error Bad request"
// @Router /products/size [post]
// @Security Bearer
func Create(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var size products.SizeType
		authHeader := ctx.GetHeader("Authorization")
		splitted := strings.Split(authHeader, " ")

		token := splitted[1]

		if err := ctx.ShouldBindJSON(&size); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		output, err := aws.NewConnect().Cognito.GetUsername(ctx, token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
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

		if err := size.Insert(db, *output.Username); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"message": "size type successfully created",
		})
	}
}

// @Summary GetAll Size process
// @Description do get all size
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {object} dtos.MessagesResponses "get all size"
// @Failure 400 {string} string "cookie not found"
// @Router /products/size [get]
// @Security Bearer
func GetAll(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var size products.SizeType
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

		result, err := size.GetAll(db)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"size_type": result,
		})
	}
}

// @Summary GetID Size process
// @Description do get all size
// @Tags products
// @Accept json
// @Produce json
// @Param id path integer true "get id size"
// @Success 200 {object} dtos.MessagesResponses "get all size"
// @Failure 400 {string} string "cookie not found"
// @Router /products/size/{id} [get]
// @Security Bearer
func GetID(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var size products.SizeType
		id, err := strconv.Atoi(ctx.Params.ByName("id"))
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
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

		result, err := size.Get(id, db)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"size_type": result,
		})
	}
}

// @Summary Update Size process
// @Description do get all size
// @Tags products
// @Accept json
// @Produce json
// @Param id path integer true "access id size type"
// @Param name body dtos.SizeType true "update size type"
// @Success 200 {object} dtos.MessagesResponses "get all size"
// @Failure 400 {string} string "cookie not found"
// @Router /products/size/{id} [put]
// @Security Bearer
func Update(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var size products.SizeType
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
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}

		if err := ctx.ShouldBindJSON(&size); err != nil {
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
			ctx.JSON(http.StatusForbidden, gin.H{"message": "permission denied"})
			return
		}

		err = size.Update(id, *output.Username, db)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("update successfully id:%d", id),
		})
	}
}

// @Summary DeleteID Size process
// @Description do delete a size
// @Tags products
// @Accept json
// @Produce json
// @Param id path integer true "delete size with id"
// @Success 200 {object} dtos.MessagesResponses "get all size"
// @Failure 400 {string} string "cookie not found"
// @Router /products/size/{id} [delete]
// @Security Bearer
func Delete(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var size products.SizeType
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
			ctx.JSON(http.StatusNotFound, gin.H{
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
			ctx.JSON(http.StatusForbidden, gin.H{"message": "permission denied"})
			return
		}

		err = size.Delete(id, *output.Username, db)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("delete successfully id:%d", id),
		})
	}
}
