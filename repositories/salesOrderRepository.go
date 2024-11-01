package repositories

import (
	"errors"
	"sales-order-golangGin/config"
	"sales-order-golangGin/models"
)

type SalesOrderRepository struct{}

func (r *SalesOrderRepository) CreateSalesOrder(order *models.SalesOrder) error {
	db := config.DB

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

	for i := range order.Items {
		order.Items[i].Id_Order = order.Id_Order
		if err := tx.Create(&order.Items[i]).Error; err != nil {
			tx.Rollback()
			return err
		}
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
	if err := db.Preload("Items").Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
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

	if err := tx.Where("id_order = ?", id).Delete(&models.ItemOrder{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
