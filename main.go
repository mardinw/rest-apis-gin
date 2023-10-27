package main

import (
	"log"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"payuoge.com/internal/api/servers"
	"payuoge.com/pkg/cache"
	"payuoge.com/pkg/database"
)

// @contact.name   API Support
// @contact.email  cs@payuoge.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println("load env successfully")
	}

	// connection for postgres
	db, err := database.Init()
	if err != nil {
		log.Println(err.Error())
		return
	} else {
		log.Println("connection pool successfully")
	}

	servers.Migrate(db)
	defer database.CloseDB()
	// connection for cache using redis
	caches, err := cache.Init()
	if err != nil {
		log.Println(err.Error())
	}

	servers.Run(db, caches)

}
