package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tim73916/org-structure-api/internal/config"
	"github.com/Tim73916/org-structure-api/internal/database"
	"github.com/Tim73916/org-structure-api/internal/models"
	"github.com/Tim73916/org-structure-api/internal/repositories"
	"github.com/Tim73916/org-structure-api/internal/services"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	cfg := &config.Config{
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "orguser",
		DBPassword: "orgpass",
		DBName:     "orgdb",
		ServerPort: "8080",
	}
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		t.Skip("Skipping test: database not available")
	}
	// Очистка
	db.Exec("TRUNCATE TABLE employees, departments RESTART IDENTITY CASCADE")
	return db
}

func TestCreateDepartmentHandler(t *testing.T) {
	db := setupTestDB(t)
	deptRepo := repositories.NewDepartmentRepo(db)
	empRepo := repositories.NewEmployeeRepo(db)
	deptService := services.NewDepartmentService(deptRepo, empRepo, db)
	deptHandler := NewDepartmentHandler(deptService)

	body := `{"name":"Test Department"}`
	req := httptest.NewRequest("POST", "/departments/", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	deptHandler.CreateDepartment(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var dept models.Department
	err := json.Unmarshal(w.Body.Bytes(), &dept)
	assert.NoError(t, err)
	assert.Equal(t, "Test Department", dept.Name)
}

func TestGetDepartmentHandler(t *testing.T) {
	db := setupTestDB(t)
	deptRepo := repositories.NewDepartmentRepo(db)
	empRepo := repositories.NewEmployeeRepo(db)
	deptService := services.NewDepartmentService(deptRepo, empRepo, db)
	deptHandler := NewDepartmentHandler(deptService)

	dept := &models.Department{Name: "IT"}
	db.Create(dept)

	req := httptest.NewRequest("GET", "/departments/1?depth=2&include_employees=true", nil)
	w := httptest.NewRecorder()

	deptHandler.GetDepartment(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Department
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "IT", response.Name)
}
