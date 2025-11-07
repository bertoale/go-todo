package controllers

import (
	"fmt"
	"rest-api/internal/dto/request"
	"rest-api/internal/models"
	"rest-api/internal/services"

	"github.com/gofiber/fiber/v2"
)

type TaskController struct {
	taskService services.TaskService
}

func NewTaskController(taskService services.TaskService) *TaskController {
	return &TaskController{
		taskService: taskService,
	}
}

func (ctrl *TaskController) CreateTask(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	var req request.TaskCreateRequest
	if err:= c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	blog, err := ctrl.taskService.CreateTask(user.ID, req.Title, req.Description)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Task created successfully",
		"task":    blog,
	})
}

func (ctrl *TaskController) DeleteTask(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	id := c.Params("id")

	var taskID uint
	if _, err := fmt.Sscanf(id, "%d", &taskID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid task ID",
		})
	}
	if err := ctrl.taskService.DeleteTask(user.ID, taskID); err != nil {
		statusCode := fiber.StatusInternalServerError
		if err.Error() == "task not found" {
			statusCode = fiber.StatusNotFound
		} else if err.Error() == "unauthorized to delete this task" {
			statusCode = fiber.StatusForbidden
		}
		return c.Status(statusCode).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "Task deleted successfully",
	})
}

func (ctrl *TaskController) GetTaskByID(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	id := c.Params("id")

	var taskID uint
	if _, err := fmt.Sscanf(id, "%d", &taskID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid task ID",
		})
	}
	task, err := ctrl.taskService.GetTasksByID(taskID)
	if err != nil {
		statusCode := fiber.StatusInternalServerError
		if err.Error() == "task not found" {
			statusCode = fiber.StatusNotFound
		}
		return c.Status(statusCode).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if task.UserID != user.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "unauthorized to access this task",
		})
	}
	return c.JSON(fiber.Map{
		"task": task,
	})
}

func (ctrl *TaskController) GetTasksByUserID(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	tasks, err := ctrl.taskService.GetTasksByUserID(user.ID)
	if err != nil {
		statusCode := fiber.StatusInternalServerError
		if err.Error() == "no tasks found for this user" {
			statusCode = fiber.StatusNotFound
		}
		return c.Status(statusCode).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"tasks": tasks,
	})
}

func (ctrl *TaskController) UpdateTask(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	id := c.Params("id")
	var taskID uint
	if _, err := fmt.Sscanf(id, "%d", &taskID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid task ID",
		})
	}
	var req request.TaskUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}
	updatedTask, err := ctrl.taskService.UpdateTask(user.ID, taskID, req.Title, req.Description, req.IsCompleted)
	if err != nil {
		statusCode := fiber.StatusInternalServerError
		if err.Error() == "task not found" {
			statusCode = fiber.StatusNotFound
		} else if err.Error() == "unauthorized to update this task" {
			statusCode = fiber.StatusForbidden
		}	 else if err.Error() == "title or description must be provided" {
			statusCode = fiber.StatusBadRequest
		}
		return c.Status(statusCode).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "Task updated successfully",
		"task":    updatedTask,
	})
}