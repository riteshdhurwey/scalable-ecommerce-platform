package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/models"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/service"
)


type CategoryHandler struct {
	service *service.CategoryService
}


func InitCategoryHandler(service *service.CategoryService) *CategoryHandler{
	return &CategoryHandler{service: service}
}


func (h *CategoryHandler) CreateCategory(ctx *gin.Context){
	type createCategoryRequest struct {
		ID int64 `json:"id"`
		Name string `json:"name"`
		Description string `json:"description"`
	}

	var category createCategoryRequest

    bindingErr := ctx.ShouldBindJSON(&category)

	if bindingErr != nil {
		ctx.JSON(http.StatusBadRequest,gin.H{"error":bindingErr})
		return
	}

	categoryID,err := h.service.CreateCategory(ctx,category.Name,category.Description)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError,gin.H{"error":err})
		return
	}

	category.ID = categoryID

	ctx.JSON(http.StatusCreated,gin.H{"added":category,"response":"category added succesfully"})


}



func (h *CategoryHandler) GetCategories(ctx *gin.Context) {
    categories, err := h.service.GetCategories(ctx.Request.Context())
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }

    if len(categories) == 0 {
        ctx.JSON(http.StatusOK, gin.H{
            "categories": []models.Category{},
            "message":    "No categories found",
        })
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "categories": categories,
    })
}


