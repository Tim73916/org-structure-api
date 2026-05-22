package testutils

import (
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupTestDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test DB: %v", err)
	}

	CleanDB(db)

	return db
}

func CleanDB(db *gorm.DB) {
	db.Exec("TRUNCATE TABLE employees, departments RESTART IDENTITY CASCADE")
}

func CreateTestDepartment(db *gorm.DB, name string, parentID *uint) uint {
	dept := struct {
		Name     string
		ParentID *uint
	}{
		Name:     name,
		ParentID: parentID,
	}

	result := db.Table("departments").Create(&dept)
	if result.Error != nil {
		panic(result.Error)
	}

	var id uint
	db.Table("departments").Select("id").Where("name = ?", name).Scan(&id)
	return id
}
