package services

import (
	"errors"

	"github.com/Tim73916/org-structure-api/internal/models"
	"github.com/Tim73916/org-structure-api/internal/repositories"
	"gorm.io/gorm"
)

type DepartmentService struct {
	deptRepo *repositories.DepartmentRepo
	empRepo  *repositories.EmployeeRepo
	db       *gorm.DB
}

func NewDepartmentService(deptRepo *repositories.DepartmentRepo, empRepo *repositories.EmployeeRepo, db *gorm.DB) *DepartmentService {
	return &DepartmentService{
		deptRepo: deptRepo,
		empRepo:  empRepo,
		db:       db,
	}
}

func (s *DepartmentService) Create(name string, parentID *uint) (*models.Department, error) {
	dept := &models.Department{
		Name:     name,
		ParentID: parentID,
	}

	if err := dept.Validate(); err != nil {
		return nil, err
	}

	if !s.deptRepo.IsNameUnique(parentID, name, 0) {
		return nil, errors.New("department name already exists in this parent")
	}

	if parentID != nil {
		var parent models.Department
		if err := s.db.First(&parent, *parentID).Error; err != nil {
			return nil, errors.New("parent department not found")
		}
	}

	if err := s.deptRepo.Create(dept); err != nil {
		return nil, err
	}

	return dept, nil
}

func (s *DepartmentService) GetByID(id uint, depth int, includeEmployees bool) (*models.Department, error) {
	if depth < 1 {
		depth = 1
	}
	if depth > 5 {
		depth = 5
	}

	dept, err := s.deptRepo.GetByID(id, depth, includeEmployees)
	if err != nil {
		return nil, errors.New("department not found")
	}

	return dept, nil
}

func (s *DepartmentService) Update(id uint, name *string, parentID *uint) (*models.Department, error) {
	dept, err := s.deptRepo.GetByID(id, 1, false)
	if err != nil {
		return nil, errors.New("department not found")
	}

	if name != nil && *name != "" {
		dept.Name = *name
		if err := dept.Validate(); err != nil {
			return nil, err
		}

		if !s.deptRepo.IsNameUnique(dept.ParentID, dept.Name, id) {
			return nil, errors.New("department name already exists in this parent")
		}
	}

	if parentID != nil {
		if cycle, _ := s.deptRepo.WouldCreateCycle(id, *parentID); cycle {
			return nil, errors.New("cannot move department to its own descendant")
		}

		if *parentID != 0 {
			var parent models.Department
			if err := s.db.First(&parent, *parentID).Error; err != nil {
				return nil, errors.New("parent department not found")
			}
			dept.ParentID = parentID
		} else {
			dept.ParentID = nil
		}
	}

	if err := s.deptRepo.Update(dept); err != nil {
		return nil, err
	}

	return dept, nil
}

func (s *DepartmentService) Delete(id uint, mode string, reassignTo *uint) error {
	_, err := s.deptRepo.GetByID(id, 1, false)
	if err != nil {
		return errors.New("department not found")
	}
	if mode == "" {
		mode = "cascade"
	}

	if mode == "reassign" {
		if reassignTo == nil {
			return errors.New("reassign_to_department_id required for reassign mode")
		}

		var target models.Department
		if err := s.db.First(&target, *reassignTo).Error; err != nil {
			return errors.New("target department not found")
		}

		if id == *reassignTo {
			return errors.New("cannot reassign to the same department")
		}

		if err := s.deptRepo.ReassignEmployees(id, *reassignTo); err != nil {
			return err
		}

		descendants, _ := s.deptRepo.GetAllDescendantIDs(id)
		for _, descendantID := range descendants {
			s.deptRepo.ReassignEmployees(descendantID, *reassignTo)
		}

		return s.deptRepo.Delete(id)
	}

	return s.deptRepo.Delete(id)
}
