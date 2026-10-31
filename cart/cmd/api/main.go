package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/cart/database"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/cart/internal/handler"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/cart/internal/middleware"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/cart/internal/repository"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/cart/internal/service"
)


func main(){

	DB,err := database.ConnectDB()
	
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected To Database")

	defer DB.Close()

	engine := gin.Default()


	repo := repository.InitCartRepo(DB)

	service := service.InitCartSevice(repo)

	handler := handler.InitCartHandler(service)

	protected := engine.Group("/api/v1/cart")
	protected.Use(middleware.AuthMiddleware)

	protected.GET("",handler.GetCart)
	protected.POST("/product", handler.AddToCart)

	engine.Run(":8082")
}