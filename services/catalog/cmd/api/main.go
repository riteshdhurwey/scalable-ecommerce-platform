package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/database"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/handler"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/repository"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/service"
)


func main(){

	DB,err := database.ConnectDB()

	if err != nil {
		panic(err)
	}

	defer DB.Close()

	fmt.Println("Connected To Database")

	router := gin.Default()

	repo := repository.NewPostgreRepo(DB)

	categoryService := service.InitCategoryService(*repo)
	categoryHandler := handler.InitCategoryHandler(categoryService)

	productService := service.InitProductService(repo)
	productHandler := handler.InitProductHandler(productService)

	router.POST("/api/v1/category",categoryHandler.CreateCategory)
	router.GET("/api/v1/categories",categoryHandler.GetCategories)

	router.POST("/api/v1/product",productHandler.CreateProduct)
	router.GET("/api/v1/product/:id",productHandler.GetProductByID)
	router.DELETE("/api/v1/product/:id",productHandler.DeleteProduct)
	router.GET("/api/v1/products",productHandler.ListProducts)
	router.PATCH("/api/v1//product/:id",productHandler.UpdateProduct)
	router.GET("/api/v1/categories/:category_id/products",productHandler.ListProductsByCategory)

	router.Run(":8080")
}