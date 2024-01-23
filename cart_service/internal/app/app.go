package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"store/cart_service/docs"
	"store/cart_service/internal/controller"
	"store/cart_service/internal/jwt"
	"store/cart_service/internal/oauth"
	"store/cart_service/internal/usecase"
	"store/cart_service/internal/usecase/api"
	"store/cart_service/internal/usecase/repo"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	// swagger embed files
)

func Run() {
	pgConnStr, ok := os.LookupEnv("CART_DB_CONNECTION_STRING")

	if !ok {
		log.Fatal(fmt.Errorf("no db connection string set up in env"))
	}

	jwkServerURL, ok := os.LookupEnv("JWK_SERVER_URL")
	if !ok {
		log.Fatal(fmt.Errorf("no jwk server url set up in env"))
	}

	productApiURL, ok := os.LookupEnv("PRODUCT_SERVICE_URL")
	if !ok {
		log.Fatal(fmt.Errorf("no jwk server url set up in env"))
	}

	keyApi, err := oauth.NewJwkWebApi(jwkServerURL)
	if err != nil {
		log.Fatal(err)
	}

	pg, err := pgxpool.Connect(context.Background(), pgConnStr)
	if err != nil {
		log.Fatal(err)
	}
	keyStorage, err := oauth.NewJwkStorage(keyApi)
	if err != nil {
		log.Fatal(err)
	}
	jwtParser := jwt.NewJwtParser(keyStorage)

	cartRepo, err := repo.NewCart(pg)
	if err != nil {
		log.Fatal(err)
	}

	productApi, err := api.NewProduct(productApiURL)
	if err != nil {
		log.Fatal(err)
	}

	cartUsecase := usecase.NewCartUsecase(cartRepo, productApi)

	//controller

	cartController := controller.NewCart(cartUsecase)
	//router
	router := gin.Default()

	api := router.Group("/api")

	v1 := api.Group("/v1")
	router.Use(oauth.Middleware(jwtParser))

	v1.Use(oauth.Middleware(jwtParser))
	{
		//cart
		cart := v1.Group("/cart")
		cart.GET("/", cartController.GetCart)
		cart.POST("/", cartController.AddOrUpdate)
		cart.DELETE("/", cartController.RemoveFromCart)
		cart.GET("/products", cartController.GetProducts)

	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	docs.SwaggerInfo.BasePath = "/api/vi"
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
