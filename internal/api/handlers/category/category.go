package category

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

// @Summary create category process
// @Description do create category
// @Tags products
// @Accept json
// @Produce json
// @Param category body dtos.CategoryProducts true "create category"
// @Success 200 {object} dtos.MessagesResponses "create category successfully"
// @Failure 400 {string} string "Error Bad request"
// @Router /products/category [post]
// @Security Bearer
func Create(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var category products.CategoryProducts

		authHeader := ctx.GetHeader("Authorization")
		splitted := strings.Split(authHeader, " ")

		token := splitted[1]

		if err := ctx.ShouldBindJSON(&category); err != nil {
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

		if err := category.Insert(db, *output.Username); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"message": fmt.Sprintf("create category %s successfully", category.Name),
		})
	}
}

// @Summary GetAll category process
// @Description do get all category
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {object} dtos.MessagesResponses "get all category"
// @Failure 400 {string} string "cookies not found"
// @Router /products/category [get]
// @Security Bearer
func GetAll(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var category products.CategoryProducts
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

		result, err := category.GetAll(db)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"category": result})
	}
}

// @Summary GetID category process
// @Description do get id category
// @Tags products
// @Accept json
// @Produce json
// @Param id path integer true "get id category"
// @Success 200 {object} dtos.MessagesResponses "get category"
// @Failure 400 {string} string "cookie not found"
// @Router /products/category/{id} [get]
// @Security Bearer
func GetID(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var category products.CategoryProducts

		id, err := strconv.Atoi(ctx.Params.ByName("id"))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
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

		result, err := category.Get(id, db)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"category": result})
	}
}

// @Summary update category process
// @Description do update category
// @Tags products
// @Accept json
// @Produce json
// @Param id path integer true "access id category"
// @Param category body dtos.CategoryProducts true "update category"
// @Success 200 {object} dtos.MessagesResponses "update category successfully"
// @Failure 400 {string} string "Error Bad request"
// @Router /products/category/{id} [put]
// @Security Bearer
func Update(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var category products.CategoryProducts

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
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		if err := ctx.ShouldBindJSON(&category); err != nil {
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

		err = category.Update(id, *output.Username, db)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("update successfully id:%d", id),
		})
	}
}

// @Summary DeleteID category process
// @Description do delete a category
// @Tags products
// @Accept json
// @Produce json
// @Param id path integer true "delete size with category"
// @Success 200 {object} dtos.MessagesResponses "get all size"
// @Failure 400 {string} string "cookie not found"
// @Router /products/category/{id} [delete]
// @Security Bearer
func Delete(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var category products.CategoryProducts

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

		err = category.Delete(id, *output.Username, db)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("delete successfully id: %d", id),
		})
	}
}
