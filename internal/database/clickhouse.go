package database

import (
	"fmt"

	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

type ClickDatabase struct {
	Instance *gorm.DB
}

// запуск кликхауса.
func (currentlDB *ClickDatabase) ClickHouseStart() {
	currentlDB.checkDatabaseCreated()
	currentlDB.cOpenConnect()
	currentlDB.cHouseMigrate()
}

func (currentlDB *ClickDatabase) checkDatabaseCreated() error {

	host := "127.0.0.1"
	port := 9000
	user := "default"
	pass := "qwerty123"

	// открытие соеднение с базой по стандарту
	dsn := fmt.Sprintf("clickhouse://%s:%s@%s:%v/default?dial_timeout=10s&read_timeout=20s", user, pass, host, port)
	db, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	db.Exec("CREATE DATABASE IF NOT EXISTS rotator")

	return nil
}

func (currentlDB *ClickDatabase) cHouseMigrate() {
	currentlDB.Instance.AutoMigrate(Statistic{})
}

func (currentlDB *ClickDatabase) cOpenConnect() error {
	host := "127.0.0.1"
	port := 9000
	user := "default"
	pass := "qwerty123"

	dsn := fmt.Sprintf("clickhouse://%s:%s@%s:%v/rotator?dial_timeout=10s&read_timeout=20s", user, pass, host, port)
	db, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	currentlDB.Instance = db
	return nil
}

func (currentlDB *ClickDatabase) ClickHouseInsert(data Statistic) error {

	if result := currentlDB.Instance.Create(&data); result.Error != nil {
		return result.Error
	}

	return nil
}
