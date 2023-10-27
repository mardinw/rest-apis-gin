package servers

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sethvargo/go-envconfig"
	"payuoge.com/configs"
	"payuoge.com/internal/api/routes"
)

func Run(db *sql.DB, caches *redis.Client) error {
	var config configs.AppConfiguration

	if err := envconfig.Process(context.Background(), &config); err != nil {
		log.Fatal(err.Error())
	}

	if config.AppEnv == "production" {
		gin.SetMode("release")
		log.Printf("server %s listening on port: %d", config.AppEnv, config.Port)
	} else {
		gin.SetMode("debug")
		log.Printf("server %s listening on port: %d", config.AppEnv, config.Port)
	}

	// running server
	srv := &http.Server{
		Addr:         ":" + strconv.Itoa(config.Port),
		Handler:      routes.NewRoutes(db, caches),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen:%s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	log.Println("Shutting down server:", sig)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown:", err)
	}
	log.Println("Server exiting")

	// end
	return nil
}
