package services

import (
	"errors"
	"time"

	"github.com/Tim73916/org-structure-api/internal/models"
	"github.com/Tim73916/org-structure-api/internal/repositories"
)

type EmployeeService struct {
	empRepo  *repositories.EmployeeRepo
	deptRepo *repositories.DepartmentRepo
}

func NewEmployeeService(empRepo *repositories.EmployeeRepo, deptRepo *repositories.DepartmentRepo) *EmployeeService {
	return &EmployeeService{
		empRepo:  empRepo,
		deptRepo: deptRepo,
	}
}

func (s *EmployeeService) Create(departmentID uint, fullName, position string, hiredAt *time.Time) (*models.Employee, error) {
	exists, err := s.empRepo.DepartmentExists(departmentID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("department not found")
	}

	emp := &models.Employee{
		DepartmentID: departmentID,
		FullName:     fullName,
		Position:     position,
		HiredAt:      hiredAt,
	}

	if err := emp.Validate(); err != nil {
		return nil, err
	}

	if err := s.empRepo.Create(emp); err != nil {
		return nil, err
	}

	return emp, nil
}

func (s *EmployeeService) GetByID(id uint) (*models.Employee, error) {
	emp, err := s.empRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("employee not found")
	}
	return emp, nil
}

func (s *EmployeeService) List(departmentID *uint) ([]models.Employee, error) {
	return s.empRepo.List(departmentID)
}

func (s *EmployeeService) Update(id uint, fullName, position *string, hiredAt *time.Time) (*models.Employee, error) {
	emp, err := s.empRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	if fullName != nil && *fullName != "" {
		emp.FullName = *fullName
	}
	if position != nil && *position != "" {
		emp.Position = *position
	}
	if hiredAt != nil {
		emp.HiredAt = hiredAt
	}

	if err := emp.Validate(); err != nil {
		return nil, err
	}

	if err := s.empRepo.Update(emp); err != nil {
		return nil, err
	}

	return emp, nil
}

func (s *EmployeeService) Delete(id uint) error {
	_, err := s.empRepo.GetByID(id)
	if err != nil {
		return errors.New("employee not found")
	}

	return s.empRepo.Delete(id)
}
