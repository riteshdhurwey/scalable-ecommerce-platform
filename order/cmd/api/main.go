package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/order/database"
	cartclient "github.com/riteshdhurwey/scalable-ecommerce-platform/order/internal/client/cart"
	inventoryclient "github.com/riteshdhurwey/scalable-ecommerce-platform/order/internal/client/inventory"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/order/internal/handler"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/order/internal/middleware"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/order/internal/repository"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/order/internal/server"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/order/internal/service"
)

func main(){
	DB,err := database.ConnectDB()

	if err != nil{
		panic(err)
	}

	defer DB.Close()

    // grpc client catalog
	catalogAddr := os.Getenv("CATALOG_GRPC_ADDR")

	if catalogAddr == ""{
		catalogAddr = "catalog:50051"
	}

	inventoryClient,err := inventoryclient.InitInventoryClient(catalogAddr)

	if err != nil {
		log.Fatalf("failed to initialize client :%v",err)
	}
	defer inventoryClient.Close()

	// grpc client cart
	cartAddr := os.Getenv("CART_GRPC_ADDR")

	if cartAddr == ""{
		cartAddr = "cart:50052"
	}

	cartClient,err := cartclient.InitCartClient(cartAddr)

	if err != nil {
		log.Fatalf("failed to initialize cart client : %v",err)
	}

	defer cartClient.Close()

	
	repo := repository.InitOrderRepo(DB)
	
	service := service.InitOrderService(repo,inventoryClient,cartClient)
	
	Orderhandler := handler.InitOrderHandler(service)
	// grpc server 

	orderGrpcHandler := handler.InitOrderGrpcHandler(service)

	go func() {
		
		err := server.StartGrpcServer(":50053", orderGrpcHandler)

		if err != nil {
			log.Fatalf("gRPC server failed %v",err)
		}
	}()
	
	router := gin.Default()
	

	protected := router.Group("/api/v1/order")
	protected.Use(middleware.AuthMiddleware)

	protected.POST("",Orderhandler.CreateOrder)

	router.Run(":8083")
}