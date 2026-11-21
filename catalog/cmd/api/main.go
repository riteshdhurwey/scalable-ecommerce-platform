package main

import (
	"context"
	"fmt"

	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/database"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/handler"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/repository"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/server"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/service"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/worker"

	orderclient "github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/clients/order"
)


func main(){

	DB,err := database.ConnectDB()

	if err != nil {
		panic(err)
	}

	defer DB.Close()

	fmt.Println("Connected To Database")

	//grpc client order 

	orderAddr := os.Getenv("ORDER_GRPC_ADDR")

	if orderAddr == ""{
		orderAddr = "order:50053"
	}

	orderClient,err := orderclient.InitOrderClient(orderAddr)

	log.Printf("Connecting to inventory service: %s", orderAddr)
	if err != nil {
	   log.Fatalf("failed to init order client %v",err)
	}

	defer orderClient.Close()

	repo := repository.NewPostgreRepo(DB)

	categoryService := service.InitCategoryService(*repo)
	categoryHandler := handler.InitCategoryHandler(categoryService)
	InventoryRepo := repository.InitInventoryRepo(DB)
	InventoryService := service.InitInventoryService(InventoryRepo)
	_= handler.InitInventoryHandler(InventoryService)

	productService := service.InitProductService(repo)
	productHandler := handler.InitProductHandler(productService)

	catalogHandler := handler.InitCatalogGRPCHandler(productService,InventoryService)

	// GRPC SERVER
	go func() {
		if err := server.StartGRPCServer(
			":50051",
			catalogHandler,
		); err != nil {
			log.Fatalf("gRPC server failed: %v", err)
		}
	}()

	//worker releases expired reservations
	ctx,cancel := context.WithCancel(context.Background())
	defer cancel()
	go worker.StartReservationWorker(ctx,DB,orderClient)

	// http server
	router := gin.Default()

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