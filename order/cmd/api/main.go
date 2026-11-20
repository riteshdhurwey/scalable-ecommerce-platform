package main

import (
	"log"
	"os"

	inventoryclient "github.com/riteshdhurwey/scalable-ecommerce-platform/order/internal/client/inventory"
	cartclient "github.com/riteshdhurwey/scalable-ecommerce-platform/order/internal/client/cart"
	"github.com/gin-gonic/gin"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/order/database"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/order/internal/handler"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/order/internal/middleware"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/order/internal/repository"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/order/internal/service"
)


func main(){
	DB,err := database.ConnectDB()

	if err != nil{
		panic(err)
	}

	defer DB.Close()

	catalogAddr := os.Getenv("CATALOG_GRPC_ADDR")

	if catalogAddr == ""{
		catalogAddr = "catalog:50051"
	}

	inventoryClient,err := inventoryclient.InitInventoryClient(catalogAddr)

	if err != nil {
		log.Fatalf("failed to initialize client :%v",err)
	}
	defer inventoryClient.Close()

	cartAddr := os.Getenv("CART_GRPC_ADDR")

	if cartAddr == ""{
		cartAddr = "cart:50052"
	}

	cartClient,err := cartclient.InitCartClient(cartAddr)

	if err != nil {
		log.Fatalf("failed to initialize cart client : %v",err)
	}

	defer cartClient.Close()

	router := gin.Default()

	repo := repository.InitOrderRepo(DB)

	service := service.InitOrderService(repo,inventoryClient,cartClient)

	handler := handler.InitOrderHandler(service)

	protected := router.Group("/api/v1/order")
	protected.Use(middleware.AuthMiddleware)

	protected.POST("",handler.CreateOrder)

	router.Run(":8083")
}