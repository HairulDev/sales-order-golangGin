package routes

import (
	"sales-order-golangGin/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, salesOrderController *controllers.SalesOrderController) {
	api := router.Group("/api")
	{
		api.POST("/salesorder", salesOrderController.CreateSalesOrder)
		api.GET("/salesorder", salesOrderController.GetSalesOrders)
		api.GET("/salesorder/:id", salesOrderController.GetSalesOrderById)
		api.PUT("/salesorder/:id", salesOrderController.UpdateSalesOrder)
		api.DELETE("/salesorder/:id", salesOrderController.DeleteSalesOrder)
	}
}
