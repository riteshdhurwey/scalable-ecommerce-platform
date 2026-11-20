package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/database"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/handler"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/repository"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/server"
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
	InventoryRepo := repository.InitInventoryRepo(DB)
	InventoryService := service.InitInventoryService(InventoryRepo)
	_= handler.InitInventoryHandler(InventoryService)

	productService := service.InitProductService(repo)
	productHandler := handler.InitProductHandler(productService)

	catalogHandler := handler.InitCatalogGRPCHandler(productService,InventoryService)

	go func() {
		if err := server.StartGRPCServer(
			":50051",
			catalogHandler,
		); err != nil {
			log.Fatalf("gRPC server failed: %v", err)
		}
	}()

	router.POST("/api/v1/category",categoryHandler.CreateCategory)
	router.GET("/api/v1/categories",categoryHandler.GetCategories)

	router.POST("/api/v1/product",productHandler.CreateProduct)
	router.GET("/api/v1/product/:id",productHandler.GetProductByID)
	router.DELETE("/api/v1/product/:id",productHandler.DeleteProduct)
	router.PATCH("/api/v1/product/:id",productHandler.UpdateProduct)
	router.GET("/api/v1/products",productHandler.ListProducts)
	router.GET("/api/v1/categories/:category_id/products",productHandler.ListProductsByCategory)
	router.POST("/api/v1/products/batch",productHandler.GetProductsByIDs)

	// router.POST("/api/v1/inventory/reserve",InventoryHandler.ReserveInventory)

	router.Run(":8080")
}