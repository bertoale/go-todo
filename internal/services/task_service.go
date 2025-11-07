package services

import (
	"errors"
	"rest-api/internal/models"
	"rest-api/internal/repositories"

	"gorm.io/gorm"
)

type TaskService interface {
	CreateTask(userID uint, title, description string) (*models.Task, error)
	GetTasksByUserID(userID uint) ([]models.Task, error)
	GetTasksByID(id uint) (*models.Task, error)
	UpdateTask(userID, blogID uint, title, description *string, isCompleted *bool) (*models.Task, error)
	DeleteTask(userID, taskID uint) error
}

type taskService struct {
	taskRepo repositories.TaskRepository
}

// CreateTask implements TaskService.
func (t *taskService) CreateTask(userID uint, title string, description string) (*models.Task, error) {
	if title == "" && description == "" {
		return nil, errors.New("title or description must be provided")
	}	
	task := &models.Task{
		UserID:      userID,
		Title:       title,
		Description: description,
	}
	if err := t.taskRepo.Create(task);  err != nil {
		return nil, errors.New("failed to create task")
	}
	return task, nil
}

// DeleteTask implements TaskService.
func (t *taskService) DeleteTask(userID uint, taskID uint) error {
	task, err := t.taskRepo.FindByID(taskID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("task not found")
		}
		return errors.New("failed to retrieve task")
	}
	if task.UserID != userID {
		return errors.New("unauthorized to delete this task")
	}
	if err := t.taskRepo.Delete(task); err != nil {
		return errors.New("failed to delete task")
	}
	return nil
}

// GetTasksByID implements TaskService.
func (t *taskService) GetTasksByID(id uint) (*models.Task, error) {
	task, err := t.taskRepo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("task not found")
		}
		return nil, errors.New("failed to retrieve task")
	}
	return task, nil
}

// GetTasksByUserID implements TaskService.
func (t *taskService) GetTasksByUserID(userID uint) ([]models.Task, error) {
	tasks, err := t.taskRepo.FindAllByUserID(userID)
	if err != nil {
		return nil, errors.New("failed to retrieve tasks")
	}

	if len(tasks) == 0 {
		return nil, errors.New("no tasks found for this user")
	}

	return tasks, nil
}


// UpdateTask implements TaskService.
func (t *taskService) UpdateTask(userID uint, taskID uint, title *string, description *string, isCompleted *bool) (*models.Task, error) {
	// 1️⃣ Ambil task berdasarkan ID
	task, err := t.taskRepo.FindByID(taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("task not found")
		}
		return nil, errors.New("failed to retrieve task")
	}

	// 2️⃣ Pastikan task milik user yang sedang login
	if task.UserID != userID {
		return nil, errors.New("unauthorized to update this task")
	}

	// 3️⃣ Update field yang dikirim (gunakan pointer agar bisa optional)
	if title != nil {
		task.Title = *title
	}
	if description != nil {
		task.Description = *description
	}
	if isCompleted != nil {
		task.IsCompleted = *isCompleted
	}

	// 4️⃣ Simpan perubahan ke database
	if err := t.taskRepo.Update(task); err != nil {
		return nil, errors.New("failed to update task")
	}

	return task, nil
}


func NewTaskService(taskRepo repositories.TaskRepository) TaskService {
	return &taskService{taskRepo: taskRepo}
}
