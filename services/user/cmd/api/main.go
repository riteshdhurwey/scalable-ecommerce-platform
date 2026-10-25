package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/user/database"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/user/internal/handler"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/user/internal/repository"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/user/internal/routes"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/user/internal/service"
)

func main() {

	DB, err := database.ConnectDB()

	if err != nil {
		panic(err)
	}

	fmt.Println("Connected To Database")

	defer DB.Close()

	engine := gin.Default()

	repo := repository.InitUserRepo(DB)

	userService := service.InitUserService(repo)

	userHandler := handler.InitUserHandler(userService)

	router := routes.InitRouter(engine, userHandler)

	router.UserRoutes()

	engine.Run(":8081")
}
