package main

import (
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

	router := gin.Default()

	repo := repository.InitOrderRepo(DB)

	service := service.InitOrderService(repo)

	handler := handler.InitOrderHandler(service)

	protected := router.Group("/api/v1/order")
	protected.Use(middleware.AuthMiddleware)

	protected.POST("",handler.CreateOrder)

	router.Run(":8083")
}