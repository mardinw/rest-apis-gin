package routes

import (
	"context"
	"database/sql"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sethvargo/go-envconfig"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"payuoge.com/configs"
	"payuoge.com/docs"
	"payuoge.com/internal/api/handlers"
	"payuoge.com/internal/api/handlers/category"
	"payuoge.com/internal/api/handlers/operationals"
	"payuoge.com/internal/api/handlers/products"
	"payuoge.com/internal/api/handlers/size"
	"payuoge.com/internal/api/handlers/transactions"
	"payuoge.com/internal/api/middleware"
)

func NewRoutes(db *sql.DB, caches *redis.Client) *gin.Engine {
	var config configs.AppConfiguration
	if err := envconfig.Process(context.Background(), &config); err != nil {
		log.Fatal(err.Error())
	}

	docs.SwaggerInfo.Title = "Pay API"
	docs.SwaggerInfo.Description = "This is Pay API Server"
	docs.SwaggerInfo.Version = "1.0.0"
	if config.AppEnv == "production" {
		docs.SwaggerInfo.Host = config.Swag.Host
	} else {
		docs.SwaggerInfo.Host = "localhost:4001"
	}
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	router := gin.Default()

	router.Use(gin.Logger())

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000", "https://play.google.com"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
	}))

	router.HandleMethodNotAllowed = true

	v1 := router.Group("/v1")
	{
		v1.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		v1.GET("/health", handlers.HealthCheck)

		// auth group
		auth := v1.Group("/auth")
		{
			auth.POST("/login", middleware.Login(caches))
			auth.POST("/register", middleware.Register)
			auth.POST("/confirm", middleware.Confirmation)
			auth.POST("/resend", middleware.ResendCode)
			auth.POST("/forgot", middleware.ForgotPassword)
			auth.POST("/reset", middleware.ResetPassword)
			auth.GET("/logout", middleware.Logout(caches))
			auth.POST("/add-user-groups", middleware.AddUserToGroup(caches))
			auth.GET("/list-groups", middleware.GetGroups())
			auth.GET("/google", middleware.LoginGA)
			auth.GET("/callback", middleware.CallbackCognito(caches))
		}

		// all product
		productAllHand := v1.Group("/products")
		productAllHand.Use(middleware.Auth(caches))
		{
			productAllHand.GET("", products.GetAll(db))

			// category
			categoryHand := productAllHand.Group("/category")
			categoryHand.Use(middleware.Auth(caches))
			{
				categoryHand.POST("", category.Create(db))
				categoryHand.GET("", category.GetAll(db))
				categoryHand.GET("/:id", category.GetID(db))
				categoryHand.PUT("/:id", category.Update(db))
				categoryHand.DELETE("/:id", category.Delete(db))
			}

			// size type
			sizeTypeHand := productAllHand.Group("/size")
			sizeTypeHand.Use(middleware.Auth(caches))
			{
				sizeTypeHand.POST("", size.Create(db))
				sizeTypeHand.GET("", size.GetAll(db))
				sizeTypeHand.GET("/:id", size.GetID(db))
				sizeTypeHand.PUT("/:id", size.Update(db))
				sizeTypeHand.DELETE("/:id", size.Delete(db))
			}
		}

		// product group
		productHand := v1.Group("/product")
		productHand.Use(middleware.Auth(caches))
		{
			productHand.GET("/:id", products.GetID(db))
		}

		// groceries group
		groceriesHand := v1.Group("/groceries")
		groceriesHand.Use(middleware.Auth(caches))
		{
			operateHand := groceriesHand.Group("/operational")
			{
				operateHand.GET("", operationals.GetAll(db))
				operateHand.POST("", operationals.Create(db))
				operateHand.PUT("/:id", operationals.Update(db))
				operateHand.GET("/:id", operationals.GetId(db))
				operateHand.DELETE("/:id", operationals.DeleteID(db))
			}
			productGroceriesHand := groceriesHand.Group("/products")
			{
				productGroceriesHand.POST("", products.Create(db))
				productGroceriesHand.GET("", products.GetProductsGrocery(db)) // nanti pakai id grosir
				productGroceriesHand.PUT("/:id", products.Update(db))
				productGroceriesHand.DELETE("/:id", products.DeleteID(db))
			}
		}

		// transaction group
		transactionHand := v1.Group("/transactions")
		transactionHand.Use(middleware.Auth(caches))
		{
			// carts
			transactionCartHand := transactionHand.Group("/carts")
			{
				transactionCartHand.POST("", transactions.CreateCart(db))
				transactionCartHand.GET("", transactions.GetAllCart(db))
				transactionCartHand.GET("/:id", transactions.GetIDCart(db))
				transactionCartHand.PUT("/:id", transactions.UpdateChart(db))
				transactionCartHand.DELETE("/:id", transactions.DeleteChart(db))
			}

			// checkout
			transactionCheckoutHand := transactionHand.Group("/checkout")
			{
				transactionCheckoutHand.POST("", transactions.CreateCheckout(db))
				transactionCheckoutHand.GET("", transactions.GetCheckout(db))
				transactionCheckoutHand.GET("/:id", transactions.GetIDCheckout(db))
				transactionCheckoutHand.PUT("/:id", transactions.UpdateCheckout(db))
				transactionCheckoutHand.DELETE("/:id", transactions.DeleteCheckout(db))
			}

			// orders
			transactionOrderHand := transactionHand.Group("/orders")
			{
				transactionOrderHand.POST("", transactions.CreateOrders(db))
				transactionOrderHand.GET("", transactions.GetOrders(db))
				transactionOrderHand.GET("/:id", transactions.GetIDOrder(db))
				transactionOrderHand.PUT("/:id", transactions.UpdateOrder(db))
				transactionOrderHand.DELETE("/:id", transactions.DeleteOrder(db))
			}
		}
	}

	return router
}
