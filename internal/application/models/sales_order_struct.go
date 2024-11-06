package models

import (
	"time"
)

type SalesOrder struct {
	Id_Order   string      `gorm:"primaryKey" json:"id_order"`
	Number_Order string    `json:"number_order"`
	Date      time.Time   `json:"date"`
	Customer  string      `json:"customer"`
	Address   string      `json:"address"`
	Items     []ItemOrder `gorm:"foreignKey:Id_Order" json:"items"`
}
