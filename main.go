package main

import (
	"sales-order-golangGin/config"
	"sales-order-golangGin/controllers"
	"sales-order-golangGin/repositories"
	"sales-order-golangGin/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	// Initialize Gin router
	r := gin.Default()

	salesOrderRepo := &repositories.SalesOrderRepository{}
	salesOrderController := controllers.NewSalesOrderController(salesOrderRepo)

	// Register routes
	routes.RegisterRoutes(r, salesOrderController)

	// Run the server
	r.Run(":8080")
}
