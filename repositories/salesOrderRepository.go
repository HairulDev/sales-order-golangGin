package repositories

import (
	"errors"
	"fmt"
	"sales-order-golangGin/config"
	"sales-order-golangGin/models"
	"time"

	"github.com/google/uuid"
)

type SalesOrderRepository struct{}

func (r *SalesOrderRepository) CreateSalesOrder(order *models.SalesOrder) error {
	db := config.DB

	order.Id_Order = uuid.New().String()

	for i := range order.Items {
		order.Items[i].Id_Item = uuid.New().String()
		order.Items[i].Id_Order = order.Id_Order
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *SalesOrderRepository) GetSalesOrders(page, limit int) ([]models.SalesOrder, int64, error) {
	db := config.DB
	var orders []models.SalesOrder
	var count int64

	if err := db.Model(&models.SalesOrder{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	if err := db.Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, count, nil
}

func (r *SalesOrderRepository) GetSalesOrderById(id string) (*models.SalesOrder, error) {
	db := config.DB
	var order models.SalesOrder

	if err := db.Preload("Items").First(&order, "id_order = ?", id).Error; err != nil {
		return nil, errors.New("sales order not found")
	}
	return &order, nil
}

func (r *SalesOrderRepository) UpdateSalesOrder(id string, updatedOrder *models.SalesOrder) error {
	db := config.DB

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&models.SalesOrder{}).Where("id_order = ?", id).Updates(updatedOrder).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("id_order = ?", id).Delete(&models.ItemOrder{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	for i := range updatedOrder.Items {
		updatedOrder.Items[i].Id_Order = id
		if err := tx.Create(&updatedOrder.Items[i]).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *SalesOrderRepository) DeleteSalesOrder(id string) error {
	db := config.DB

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Where("id_order = ?", id).Delete(&models.SalesOrder{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("id_order = ?", id).Delete(&models.ItemOrder{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// Method to count the total filtered sales orders
func (r *SalesOrderRepository) GetSearchSalesOrderCount(keywords string, date *time.Time) (int, error) {
	db := config.DB
	var count int64
	fmt.Println("Keywords:", keywords)
	fmt.Println("date:", date)

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
	db := config.DB
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
