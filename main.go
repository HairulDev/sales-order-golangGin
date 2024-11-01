package main

import (
	"sales-order-golangGin/config"
	"sales-order-golangGin/controllers"
	"sales-order-golangGin/repositories"
	"sales-order-golangGin/routes"

    "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	// Initialize Gin router
	r := gin.Default()

    // Set up CORS middleware
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:5173"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        AllowCredentials: true,
    }))

	salesOrderRepo := &repositories.SalesOrderRepository{}
	salesOrderController := controllers.NewSalesOrderController(salesOrderRepo)

	// Register routes
	routes.RegisterRoutes(r, salesOrderController)

	// Run the server
	r.Run(":8080")
}
