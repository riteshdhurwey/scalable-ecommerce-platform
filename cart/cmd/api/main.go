package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/cart/database"
	catalogclient "github.com/riteshdhurwey/scalable-ecommerce-platform/services/cart/internal/clients/catalog"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/cart/internal/handler"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/cart/internal/middleware"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/cart/internal/repository"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/cart/internal/server"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/cart/internal/service"
)


func main(){

	DB,err := database.ConnectDB()
	
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected To Database")

	defer DB.Close()

	// GRPC Client To Get Products Information From catalog
	catalogAddr := os.Getenv("CATALOG_GRPC_ADDR")
	if catalogAddr == "" {
		catalogAddr = "catalog:50051"
	}

	catalogClient, err := catalogclient.InitCatalogClient(catalogAddr)
	if err != nil {
		log.Fatalf("failed to initialize catalog client: %v", err)
	}
	defer catalogClient.Close()
	
	repo := repository.InitCartRepo(DB)
	
	service := service.InitCartSevice(repo,*catalogClient)
	
	cartHandler := handler.InitCartHandler(service)

	// GRPC Server to Serve Cart To Order Service
	cartGrpcHandler := handler.InitCartGrpcHandler(service)
	go func () {
		if err := server.StartGrpcServer(":50052",cartGrpcHandler,); err != nil {
			log.Fatalf("grpc server failed %v",err)
		}
	}()

	
	engine := gin.Default()
	

	protected := engine.Group("/api/v1/cart")
	protected.Use(middleware.AuthMiddleware)

	protected.GET("",cartHandler.CheckOut)
	protected.POST("/product", cartHandler.AddToCart)
	protected.PATCH("/product/update",cartHandler.UpdateCart)
	protected.DELETE("/product/:id",cartHandler.DeleteItem)
	engine.Run(":8082")
}