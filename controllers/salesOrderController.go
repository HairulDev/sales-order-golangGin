package controllers

import (
	"net/http"
	"sales-order-golangGin/models"
	"sales-order-golangGin/repositories"

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
	if err := c.repo.CreateSalesOrder(&order); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, order)
}

func (c *SalesOrderController) GetSalesOrders(ctx *gin.Context) {
	page := 1
	limit := 10 

	if p := ctx.Query("page"); p != "" {
		// Parse page query parameter
		// (Implement parsing and error handling)
	}
	if l := ctx.Query("limit"); l != "" {
		// Parse limit query parameter
		// (Implement parsing and error handling)
	}

	orders, count, err := c.repo.GetSalesOrders(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Sales orders retrieved successfully","status": true,"data": orders,"total": count})
}

func (c *SalesOrderController) GetSalesOrderById(ctx *gin.Context) {
	id := ctx.Param("id")
	order, err := c.repo.GetSalesOrderById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Sales order retrieved successfully","status": true,"data": order})
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