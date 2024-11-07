package repositories

import (
	"errors"
	"fmt"
	"sales-order-golangGin/internal/application/models"
	"sales-order-golangGin/internal/pkg/database/sql/configs"
	"sync"
	"time"

	"github.com/google/uuid"
)

type SalesOrderRepository struct{}

func (r *SalesOrderRepository) CreateSalesOrder(order *models.SalesOrder) error {
	db := configs.DB
	order.Id_Order = uuid.New().String()

	var wg sync.WaitGroup
	errCh := make(chan error, 1) // Only need one error channel to catch the first error

	// Start transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := tx.Create(&order).Error; err != nil {
			errCh <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := range order.Items {
			order.Items[i].Id_Item = uuid.New().String()
			order.Items[i].Id_Order = order.Id_Order

			// Insert item with error handling
			if err := tx.Create(&order.Items[i]).Error; err != nil {
				// Send error and stop further processing
				errCh <- err
				return
			}
		}
	}()

	// Wait for goroutines and check for errors
	go func() {
		wg.Wait()
		close(errCh) // Close error channel once all goroutines finish
	}()

	// Check for errors and roll back if any error was encountered
	if err, ok := <-errCh; ok {
		tx.Rollback()
		return err
	}

	// Commit transaction if no errors
	return tx.Commit().Error
}

func (r *SalesOrderRepository) GetSalesOrders(page, limit int) ([]models.SalesOrder, int, error) {
	var orders []models.SalesOrder
	var totalCount int64 // Define totalCount as int64 to be compatible with GORM's Count method
	var finalCount int   // Define finalCount as int to be used in response

	ordersCh := make(chan []models.SalesOrder)
	countCh := make(chan int)
	errCh := make(chan error, 2)

	go func() {
		err := configs.DB.Offset((page - 1) * limit).Limit(limit).Find(&orders).Error
		if err != nil {
			errCh <- err
			return
		}
		ordersCh <- orders
	}()

	go func() {
		err := configs.DB.Model(&models.SalesOrder{}).Count(&totalCount).Error
		if err != nil {
			errCh <- err
			return
		}
		finalCount = int(totalCount) // Convert int64 to int for the response
		countCh <- finalCount
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return nil, 0, err
		}
	case orders := <-ordersCh:
		return orders, <-countCh, nil
	case count := <-countCh:
		return <-ordersCh, count, nil
	case <-time.After(5 * time.Second):
		return nil, 0, errors.New("Request timed out")
	}

	return orders, finalCount, nil
}

// Method to count the total filtered sales orders
func (r *SalesOrderRepository) GetSearchSalesOrderCount(keywords string, date *time.Time) (int, error) {
	db := configs.DB
	var count int64

	query := db.Model(&models.SalesOrder{})
	if keywords != "" {
		query = query.Where("number_order LIKE ? OR customer LIKE ?", "%"+keywords+"%", "%"+keywords+"%")
	}
	if date != nil {
		query = query.Where("DATE(date) = ?", date.Format("2006-01-02"))
	}

	err := query.Count(&count).Error
	return int(count), err
}

// Method to fetch filtered sales orders
func (r *SalesOrderRepository) SearchSalesOrders(keywords string, date *time.Time, page int, limit int) ([]models.SalesOrder, error) {
	db := configs.DB
	var orders []models.SalesOrder

	query := db.Model(&models.SalesOrder{}).Order("date DESC").Offset((page - 1) * limit).Limit(limit)
	if keywords != "" {
		query = query.Where("number_order LIKE ? OR customer LIKE ?", "%"+keywords+"%", "%"+keywords+"%")
	}
	if date != nil {
		query = query.Where("DATE(date) = ?", date.Format("2006-01-02"))
	}

	err := query.Find(&orders).Error
	return orders, err
}

func (r *SalesOrderRepository) GetSalesOrderById(id string) (*models.SalesOrder, error) {
	var order models.SalesOrder

	// Use Preload to fetch associated items with the sales order
	if err := configs.DB.Preload("Items").First(&order, "id_order = ?", id).Error; err != nil {
		return nil, errors.New("sales order not found")
	}

	return &order, nil
}

func (r *SalesOrderRepository) UpdateSalesOrder(order *models.SalesOrder) error {
	errCh := make(chan error, 1)

	go func() {
		if err := configs.DB.Save(order).Error; err != nil {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case err := <-errCh:
		return err
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout while updating SalesOrder")
	}
}

func (r *SalesOrderRepository) DeleteSalesOrder(id string) error {
	errCh := make(chan error, 1)

	go func() {
		if err := configs.DB.Where("id_order = ?", id).Delete(&models.SalesOrder{}).Error; err != nil {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case err := <-errCh:
		return err
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout while deleting SalesOrder")
	}
}
