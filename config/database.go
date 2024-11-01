// config/database.go
package config

import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "log"
)

var DB *gorm.DB

func ConnectDB() {
    dsn := "root:mysqlpw@tcp(127.0.0.1:55004)/latihan?charset=utf8mb4&parseTime=True&loc=Local"
    database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database!", err)
    }

    DB = database
}
