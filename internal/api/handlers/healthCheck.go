package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-envconfig"
	"payuoge.com/configs"
)

func HealthCheck(ctx *gin.Context) {
	var config configs.AppConfiguration

	if err := envconfig.Process(ctx, &config); err != nil {
		log.Fatal(err.Error())
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "available",
		"version": config.Version,
	})

}
