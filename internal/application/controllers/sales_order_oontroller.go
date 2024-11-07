package controllers

import (
	"net/http"
	"strconv"
	"time"

	"sales-order-golangGin/internal/application/models"
	"sales-order-golangGin/internal/application/repositories"

	"github.com/gin-gonic/gin"
)

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

	errCh := make(chan error, 1)
	go func() {
		errCh <- c.repo.CreateSalesOrder(&order)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusCreated, order)
	case <-time.After(5 * time.Second):
		ctx.JSON(http.StatusGatewayTimeout, gin.H{"error": "Request timed out"})
	}
}

func (c *SalesOrderController) GetSalesOrders(ctx *gin.Context) {
	page := 1
	limit := 5

	if p := ctx.Query("page"); p != "" {
		if parsedPage, err := strconv.Atoi(p); err == nil && parsedPage > 0 {
			page = parsedPage
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
			return
		}
	}

	if l := ctx.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
			return
		}
	}

	ordersCh := make(chan []models.SalesOrder)
	countCh := make(chan int)
	errCh := make(chan error, 2)

	go func() {
		orders, count, err := c.repo.GetSalesOrders(page, limit)
		if err != nil {
			errCh <- err
			return
		}
		ordersCh <- orders
		countCh <- count
	}()

	var orders []models.SalesOrder

	select {
	case orders = <-ordersCh:
		count := <-countCh
		totalPage := (count + limit - 1) / limit
		ctx.JSON(http.StatusOK, gin.H{
			"message":     "Sales orders retrieved successfully",
			"status":      true,
			"data":        orders,
			"currentPage": page,
			"totalPage":   totalPage,
			"limit":       limit,
			"count":       count,
		})
	case err := <-errCh:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	case <-time.After(5 * time.Second):
		ctx.JSON(http.StatusGatewayTimeout, gin.H{"error": "Request timed out"})
	}
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

	// Create channels for concurrent execution
	salesOrdersCh := make(chan []models.SalesOrder)
	totalCountCh := make(chan int)
	errCh := make(chan error, 2)

	// Fetch filtered sales orders concurrently
	go func() {
		salesOrders, err := c.repo.SearchSalesOrders(keywords, date, page, limit)
		if err != nil {
			errCh <- err
			return
		}
		salesOrdersCh <- salesOrders
	}()

	// Fetch total count concurrently
	go func() {
		totalCount, err := c.repo.GetSearchSalesOrderCount(keywords, date)
		if err != nil {
			errCh <- err
			return
		}
		totalCountCh <- totalCount
	}()

	// Collect results
	var salesOrders []models.SalesOrder
	var totalCount int

	select {
	case salesOrders = <-salesOrdersCh:
	case err := <-errCh:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	select {
	case totalCount = <-totalCountCh:
	case err := <-errCh:
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

func (c *SalesOrderController) GetSearchSalesOrderCount(ctx *gin.Context) {
	// Parse query parameters
	keywords := ctx.Query("keywords")
	dateStr := ctx.Query("date")

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

	// Create a channel for concurrent execution
	countCh := make(chan int)
	errCh := make(chan error, 1)

	// Run the repository method in a goroutine
	go func() {
		totalCount, err := c.repo.GetSearchSalesOrderCount(keywords, date)
		if err != nil {
			errCh <- err
			return
		}
		countCh <- totalCount
	}()

	// Collect results
	select {
	case totalCount := <-countCh:
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Sales order count retrieved successfully",
			"status":  true,
			"count":   totalCount,
		})
	case err := <-errCh:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	case <-time.After(5 * time.Second):
		ctx.JSON(http.StatusGatewayTimeout, gin.H{"error": "Request timed out"})
	}
}

func (c *SalesOrderController) GetSalesOrderById(ctx *gin.Context) {
	id := ctx.Param("id")

	order, err := c.repo.GetSalesOrderById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Sales order retrieved successfully",
		"status":  true,
		"data":    order,
	})
}

func (c *SalesOrderController) UpdateSalesOrder(ctx *gin.Context) {
	id := ctx.Param("id")
	var order models.SalesOrder
	if err := ctx.ShouldBindJSON(&order); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order.Id_Order = id
	errCh := make(chan error, 1)
	go func() {
		errCh <- c.repo.UpdateSalesOrder(&order)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Sales order updated successfully"})
	case <-time.After(5 * time.Second):
		ctx.JSON(http.StatusGatewayTimeout, gin.H{"error": "Request timed out"})
	}
}

func (c *SalesOrderController) DeleteSalesOrder(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Sales order ID is required"})
		return
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- c.repo.DeleteSalesOrder(id)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Sales order deleted successfully"})
	case <-time.After(5 * time.Second):
		ctx.JSON(http.StatusGatewayTimeout, gin.H{"error": "Request timed out"})
	}
}
