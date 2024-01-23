package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"product_service/docs"
	"product_service/internal/controller"
	"product_service/internal/jwt"
	"product_service/internal/oauth"
	"product_service/internal/usecase"
	"product_service/internal/usecase/repo"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	// swagger embed files
)

func Run() {
	pgConnStr, ok := os.LookupEnv("DB_CONNECTION_STRING")

	if !ok {
		log.Fatal(fmt.Errorf("no db connection string set up in env"))
	}

	jwkServerURL, ok := os.LookupEnv("JWK_SERVER_URL")
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

	//repo
	productRepo, err := repo.NewProduct(pg)
	if err != nil {
		log.Fatal(err)
	}

	sellerRepo, err := repo.NewSeller(pg)
	if err != nil {
		log.Fatal(err)
	}
	//service

	productUsecase := usecase.NewProduct(productRepo)
	sellerUsecase := usecase.NewSeller(sellerRepo, productUsecase)

	//controller

	seller := controller.NewSeller(sellerUsecase)
	product := controller.NewProduct(productUsecase)
	//router
	router := gin.Default()

	api := router.Group("/api")

	v1 := api.Group("/v1")
	router.Use(oauth.Middleware(jwtParser))

	v1.Use(oauth.Middleware(jwtParser))
	{
		//product
		prod := v1.Group("/products")
		prod.GET("/:id", product.GetById)
		prod.GET("/", product.Page)
		prod.PUT("/:id", product.Update)

		//seller
		sell := v1.Group("/sellers")
		sell.GET("/", seller.GetAll)
		sell.GET("/:id", seller.GetById)
		sell.POST("/", seller.Create)
		sell.PUT("/:id", seller.Update)
		sell.GET("/:id/products", seller.GetProductPage)
		sell.POST("/:id/products", seller.AddProduct)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	docs.SwaggerInfo.BasePath = "/api/vi"
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
