package repositories

import (
	"github.com/Tim73916/org-structure-api/internal/models"
	"gorm.io/gorm"
)

type EmployeeRepo struct {
	db *gorm.DB
}

func NewEmployeeRepo(db *gorm.DB) *EmployeeRepo {
	return &EmployeeRepo{db: db}
}

func (r *EmployeeRepo) Create(emp *models.Employee) error {
	return r.db.Create(emp).Error
}

func (r *EmployeeRepo) GetByID(id uint) (*models.Employee, error) {
	var emp models.Employee
	err := r.db.First(&emp, id).Error
	if err != nil {
		return nil, err
	}
	return &emp, nil
}

func (r *EmployeeRepo) List(departmentID *uint) ([]models.Employee, error) {
	var employees []models.Employee
	query := r.db.Order("created_at ASC")

	if departmentID != nil {
		query = query.Where("department_id = ?", *departmentID)
	}

	err := query.Find(&employees).Error
	return employees, err
}

func (r *EmployeeRepo) Update(emp *models.Employee) error {
	return r.db.Save(emp).Error
}

func (r *EmployeeRepo) Delete(id uint) error {
	return r.db.Delete(&models.Employee{}, id).Error
}

func (r *EmployeeRepo) DepartmentExists(departmentID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Department{}).Where("id = ?", departmentID).Count(&count).Error
	return count > 0, err
}

func (r *EmployeeRepo) GetByDepartmentID(departmentID uint) ([]models.Employee, error) {
	var employees []models.Employee
	err := r.db.Where("department_id = ?", departmentID).Find(&employees).Error
	return employees, err
}
