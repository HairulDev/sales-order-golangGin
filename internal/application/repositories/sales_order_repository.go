package repositories

import (
	"errors"
	"sales-order-golangGin/internal/application/models"
	"sales-order-golangGin/internal/pkg/database/sql/configs"
	"time"

	"github.com/google/uuid"
)

type SalesOrderRepository struct{}

func (r *SalesOrderRepository) CreateSalesOrder(order *models.SalesOrder) error {
	db := configs.DB
	order.Id_Order = uuid.New().String()

	for i := range order.Items {
		order.Items[i].Id_Item = uuid.New().String()
		order.Items[i].Id_Order = order.Id_Order
	}

	errChan := make(chan error)
	go func() {
		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		if err := tx.Create(&order).Error; err != nil {
			tx.Rollback()
			errChan <- err
			return
		}

		errChan <- tx.Commit().Error
	}()

	return <-errChan
}

func (r *SalesOrderRepository) GetSalesOrders(page, limit int) ([]models.SalesOrder, int64, error) {
	db := configs.DB
	var orders []models.SalesOrder
	var count int64

	countChan := make(chan error)
	ordersChan := make(chan error)

	go func() {
		countChan <- db.Model(&models.SalesOrder{}).Count(&count).Error
	}()

	go func() {
		offset := (page - 1) * limit
		ordersChan <- db.Offset(offset).Limit(limit).Find(&orders).Error
	}()

	if err := <-countChan; err != nil {
		return nil, 0, err
	}
	if err := <-ordersChan; err != nil {
		return nil, 0, err
	}

	return orders, count, nil
}

func (r *SalesOrderRepository) GetSalesOrderById(id string) (*models.SalesOrder, error) {
	db := configs.DB
	var order models.SalesOrder
	orderChan := make(chan error)

	go func() {
		orderChan <- db.Preload("Items").First(&order, "id_order = ?", id).Error
	}()

	if err := <-orderChan; err != nil {
		return nil, errors.New("sales order not found")
	}

	return &order, nil
}

func (r *SalesOrderRepository) UpdateSalesOrder(id string, updatedOrder *models.SalesOrder) error {
	db := configs.DB
	errChan := make(chan error)

	go func() {
		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		if err := tx.Model(&models.SalesOrder{}).Where("id_order = ?", id).Updates(updatedOrder).Error; err != nil {
			tx.Rollback()
			errChan <- err
			return
		}

		if err := tx.Where("id_order = ?", id).Delete(&models.ItemOrder{}).Error; err != nil {
			tx.Rollback()
			errChan <- err
			return
		}

		for i := range updatedOrder.Items {
			updatedOrder.Items[i].Id_Item = uuid.New().String()
			updatedOrder.Items[i].Id_Order = id
			if err := tx.Create(&updatedOrder.Items[i]).Error; err != nil {
				tx.Rollback()
				errChan <- err
				return
			}
		}

		errChan <- tx.Commit().Error
	}()

	return <-errChan
}

func (r *SalesOrderRepository) DeleteSalesOrder(id string) error {
	db := configs.DB
	errChan := make(chan error)

	go func() {
		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		if err := tx.Where("id_order = ?", id).Delete(&models.SalesOrder{}).Error; err != nil {
			tx.Rollback()
			errChan <- err
			return
		}

		if err := tx.Where("id_order = ?", id).Delete(&models.ItemOrder{}).Error; err != nil {
			tx.Rollback()
			errChan <- err
			return
		}

		errChan <- tx.Commit().Error
	}()

	return <-errChan
}

func (r *SalesOrderRepository) GetSearchSalesOrderCount(keywords string, date *time.Time) (int, error) {
	db := configs.DB
	var count int64
	countChan := make(chan error)

	go func() {
		query := db.Model(&models.SalesOrder{})
		if keywords != "" {
			query = query.Where("number_order LIKE ? OR customer LIKE ?", "%"+keywords+"%", "%"+keywords+"%")
		}
		if date != nil {
			query = query.Where("DATE(date) = ?", date.Format("2006-01-02"))
		}

		countChan <- query.Count(&count).Error
	}()

	if err := <-countChan; err != nil {
		return 0, err
	}

	return int(count), nil
}

func (r *SalesOrderRepository) SearchSalesOrders(keywords string, date *time.Time, page int, limit int) ([]models.SalesOrder, error) {
	db := configs.DB
	var orders []models.SalesOrder
	ordersChan := make(chan error)

	go func() {
		query := db.Model(&models.SalesOrder{}).Order("date DESC").Offset((page - 1) * limit).Limit(limit)
		if keywords != "" {
			query = query.Where("number_order LIKE ? OR customer LIKE ?", "%"+keywords+"%", "%"+keywords+"%")
		}
		if date != nil {
			query = query.Where("DATE(date) = ?", date.Format("2006-01-02"))
		}

		ordersChan <- query.Find(&orders).Error
	}()

	if err := <-ordersChan; err != nil {
		return nil, err
	}

	return orders, nil
}
