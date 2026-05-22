package repositories

import (
	"errors"

	"github.com/Tim73916/org-structure-api/internal/models"
	"gorm.io/gorm"
)

type DepartmentRepo struct {
	db *gorm.DB
}

func NewDepartmentRepo(db *gorm.DB) *DepartmentRepo {
	return &DepartmentRepo{db: db}
}

func (r *DepartmentRepo) Create(dept *models.Department) error {
	return r.db.Create(dept).Error
}

func (r *DepartmentRepo) GetByID(id uint, depth int, includeEmployees bool) (*models.Department, error) {
	var dept models.Department

	query := r.db
	if includeEmployees {
		query = query.Preload("Employees", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		})
	}

	if depth > 0 {
		preloadStr := "Children"
		for i := 1; i < depth; i++ {
			preloadStr += ".Children"
			query = query.Preload(preloadStr)
		}
	}

	err := query.First(&dept, id).Error
	if err != nil {
		return nil, err
	}
	return &dept, nil
}

func (r *DepartmentRepo) Update(dept *models.Department) error {
	return r.db.Save(dept).Error
}

func (r *DepartmentRepo) Delete(id uint) error {
	return r.db.Delete(&models.Department{}, id).Error
}

func (r *DepartmentRepo) IsNameUnique(parentID *uint, name string, excludeID uint) bool {
	var count int64
	query := r.db.Model(&models.Department{}).Where("name = ? AND id != ?", name, excludeID)

	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}

	query.Count(&count)
	return count == 0
}

func (r *DepartmentRepo) WouldCreateCycle(deptID, newParentID uint) (bool, error) {
	if deptID == newParentID {
		return true, nil
	}

	currentID := newParentID
	for {
		var parent models.Department
		err := r.db.Select("parent_id").First(&parent, currentID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				break
			}
			return false, err
		}

		if parent.ParentID != nil && *parent.ParentID == deptID {
			return true, nil
		}

		if parent.ParentID == nil {
			break
		}
		currentID = *parent.ParentID
	}
	return false, nil
}

func (r *DepartmentRepo) GetAllDescendantIDs(id uint) ([]uint, error) {
	var ids []uint
	err := r.db.Raw(`
        WITH RECURSIVE descendants AS (
            SELECT id FROM departments WHERE id = ?
            UNION ALL
            SELECT d.id FROM departments d
            INNER JOIN descendants ON descendants.id = d.parent_id
        )
        SELECT id FROM descendants WHERE id != ?
    `, id, id).Scan(&ids).Error
	return ids, err
}

func (r *DepartmentRepo) ReassignEmployees(fromDeptID, toDeptID uint) error {
	return r.db.Model(&models.Employee{}).
		Where("department_id = ?", fromDeptID).
		Update("department_id", toDeptID).Error
}
