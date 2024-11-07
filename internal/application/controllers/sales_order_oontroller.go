package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"sales-order-golangGin/internal/application/models"
	"sales-order-golangGin/internal/application/repositories"
)

// Paradigma Prosedural
type SalesOrderController struct {
	repo *repositories.SalesOrderRepository
}

func NewSalesOrderController(repo *repositories.SalesOrderRepository) *SalesOrderController {
	return &SalesOrderController{repo: repo}
}

func (c *SalesOrderController) CreateSalesOrder(ctx *gin.Context) {
	var order models.SalesOrder
	if err := ctx.ShouldBindJSON(&order); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.repo.CreateSalesOrder(&order); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, order)
}

func (c *SalesOrderController) GetSalesOrders(ctx *gin.Context) {
	// Set nilai default untuk page dan limit
	page := 1
	limit := 5

	// Ambil dan parsing query parameter 'page'
	if p := ctx.Query("page"); p != "" {
		if parsedPage, err := strconv.Atoi(p); err == nil && parsedPage > 0 {
			page = parsedPage
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
			return
		}
	}

	// Ambil dan parsing query parameter 'limit'
	if l := ctx.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
			return
		}
	}

	orders, count, err := c.repo.GetSalesOrders(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPage := int((count + int64(limit) - 1) / int64(limit))

	ctx.JSON(http.StatusOK, gin.H{
		"message":     "Sales orders retrieved successfully",
		"status":      true,
		"data":        orders,
		"currentPage": page,
		"totalPage":   totalPage,
		"limit":       limit,
		"count":       count,
	})
}

func (c *SalesOrderController) GetSalesOrderById(ctx *gin.Context) {
	id := ctx.Param("id")
	order, err := c.repo.GetSalesOrderById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Sales order retrieved successfully", "status": true, "data": order})
}

func (c *SalesOrderController) UpdateSalesOrder(ctx *gin.Context) {
	id := ctx.Param("id")
	var updatedOrder models.SalesOrder
	if err := ctx.ShouldBindJSON(&updatedOrder); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.repo.UpdateSalesOrder(id, &updatedOrder); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, updatedOrder)
}

func (c *SalesOrderController) DeleteSalesOrder(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.repo.DeleteSalesOrder(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (c *SalesOrderController) SearchSalesOrders(ctx *gin.Context) {
	// Parse query parameters
	keywords := ctx.Query("keywords")
	dateStr := ctx.Query("date")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "5"))

	// Parse date if provided
	var date *time.Time
	if dateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			return
		}
		date = &parsedDate
	}

	// Fetch filtered sales orders and total count based on search criteria
	salesOrders, err := c.repo.SearchSalesOrders(keywords, date, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalCount, err := c.repo.GetSearchSalesOrderCount(keywords, date)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPage := (totalCount + limit - 1) / limit

	ctx.JSON(http.StatusOK, gin.H{
		"message":     "Sales orders retrieved successfully",
		"status":      true,
		"data":        salesOrders,
		"currentPage": page,
		"totalPage":   totalPage,
		"limit":       limit,
		"count":       totalCount,
	})
}
