package repositories

import (
	"rest-api/internal/models"

	"gorm.io/gorm"
)

type TaskRepository interface {
	Create(task *models.Task) error
	Update(task *models.Task) error
	FindByID(id uint) (*models.Task, error)
	Delete(task *models.Task) error
	FindAllByUserID(userID uint) ([]models.Task, error)
}

type taskRepository struct {
	db *gorm.DB
}

// Create implements TaskRepository.
func (t *taskRepository) Create(task *models.Task) error {
	if err := t.db.Create(task).Error; err != nil {
		return err
	}

	// Setelah berhasil insert, ambil ulang data lengkap dengan relasi User
	if err := t.db.Preload("User").First(task, task.ID).Error; err != nil {
		return err
	}

	return nil
}


// Delete implements TaskRepository.
func (t *taskRepository) Delete(task *models.Task) error {
	return t.db.Delete(task).Error
}

// FindAllByUserID implements TaskRepository.
func (t *taskRepository) FindAllByUserID(userID uint) ([]models.Task, error) {
	var tasks []models.Task
	if err := t.db.
		Preload("User").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}


// FindByID implements TaskRepository.
func (t *taskRepository) FindByID(id uint) (*models.Task, error) {
	var tasks models.Task
	if err := t.db.Preload("User").First(&tasks, id).Error; err != nil {
		return nil, err
	}
	return &tasks, nil
}

// Update implements TaskRepository.
func (t *taskRepository) Update(task *models.Task) error {
	return t.db.Save(task).Error
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}
