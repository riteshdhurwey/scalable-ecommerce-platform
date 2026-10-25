package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/dto"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/catalog/internal/service"
)


type ProductHandler struct {
	service *service.ProductService
}

func InitProductHandler(service *service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) CreateProduct (ctx *gin.Context) {
	var request dto.ProductRequestDTO

	err := ctx.ShouldBindJSON(&request)

	if err != nil {
		ctx.JSON(http.StatusBadRequest,gin.H{"error": "Invalid query parameters"})
		return
	}

	product,err := h.service.CreateProduct(ctx,request)

	if err != nil {
    switch {
		case errors.Is(err, service.ErrInvalidProduct):
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		case errors.Is(err, service.ErrInvalidPrice):
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		case errors.Is(err, service.ErrInvalidStock):
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		case errors.Is(err, service.ErrInvalidCategory):
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		case errors.Is(err, service.ErrCategoryNotFound):
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": err.Error(),
			})

		case errors.Is(err, service.ErrProductNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})

		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":err.Error(),
			})
		}
    
		return
	}

	ctx.JSON(http.StatusCreated, product)
}


func (h *ProductHandler) GetProductByID(ctx *gin.Context){
	idstr := ctx.Param("id")

	productID , err := strconv.Atoi(idstr)

	if err != nil {
		ctx.JSON(http.StatusBadRequest,gin.H{"error":"invalid product id"})
		return
	}

	product,err := h.service.GetProductByID(ctx,int64(productID))

	if err != nil {
		if errors.Is(err,pgx.ErrNoRows){
			ctx.JSON(http.StatusNotFound,service.ErrProductNotFound.Error())
			return
			
		}

		ctx.JSON(http.StatusInternalServerError,gin.H{"error":"internal server error"})

		return
	}

	ctx.JSON(http.StatusOK,product)


}

func (h *ProductHandler) DeleteProduct(ctx *gin.Context){
	idstr := ctx.Param("id")

	productID,err := strconv.Atoi(idstr)

	if err != nil {
		ctx.JSON(http.StatusBadRequest,gin.H{"error": "Invalid query parameters"})
		return
	}

	_,deletedErr := h.service.DeleteProduct(ctx,int64(productID))

	if deletedErr != nil {
		if errors.Is(deletedErr,service.ErrProductNotFound){
			ctx.JSON(http.StatusNotFound,deletedErr.Error())
			return
		}
		ctx.JSON(http.StatusInternalServerError,gin.H{"error":"internal server error"})
		return
	}

	ctx.JSON(http.StatusNoContent,gin.H{"response":"product deleted"})

	
}


func (h *ProductHandler) UpdateProduct(ctx *gin.Context){
	var request dto.UpdateProductDTO

	err := ctx.ShouldBindJSON(&request)

	if err != nil {
		ctx.JSON(http.StatusBadRequest,gin.H{"error": "Invalid query parameters"})
		return
	}

	product,err := h.service.UpdateProduct(ctx,request)

	if err != nil {
		switch {
		case errors.Is(err,service.ErrProductNotFound) :
			ctx.JSON(http.StatusNotFound,gin.H{"error":err.Error()})
			
		case errors.Is(err,service.ErrInvalidInput) :
			ctx.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			
		default :
			ctx.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusCreated,product)
}


func (h *ProductHandler) ListProducts(ctx *gin.Context){
	pageStr := ctx.DefaultQuery("page","1")
	limitStr := ctx.DefaultQuery("limit","10")

	page,err := strconv.Atoi(pageStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest,gin.H{"error":"invalid page no"})
		return
	}

	limit , err := strconv.Atoi(limitStr)

	if err != nil {
		ctx.JSON(http.StatusBadRequest,gin.H{"error":"invalid limit"})
		return
	}

	products,err := h.service.ListProducts(ctx,page,limit)

	if err != nil {
		if errors.Is(err,service.ErrInvalidPageSize){
			ctx.JSON(http.StatusBadRequest,err.Error())
			return
		}
		ctx.JSON(http.StatusInternalServerError,err.Error())
		return
	}


	ctx.JSON(http.StatusOK,products)

}



func (h *ProductHandler) ListProductsByCategory(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")
	categoryIDStr := ctx.Param("category_id")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid page number"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1  {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "limit must be between 1 and 100"})
		return
	}

	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil || categoryID < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	products, categoryName,productErr := h.service.ListProductsByCategory(
		ctx, categoryID, page, limit,
	)
	if productErr != nil {
		if errors.Is(productErr, service.ErrCategoryNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": service.ErrCategoryNotFound.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, productErr.Error())
		return
	}

	response := dto.ProductListResponse{Page: page,Limit: limit,CategoryID: categoryID,CategoryName: categoryName,Products: products}
	

	ctx.JSON(http.StatusOK, response)
}


