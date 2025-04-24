package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dimasbagussusilo/go-clean-boilerplate/domain/entity"
	"github.com/dimasbagussusilo/go-clean-boilerplate/usecase"
)

// TaskHandler represents the HTTP handler for task operations
type TaskHandler struct {
	taskUseCase *usecase.TaskUseCase
}

// NewTaskHandler creates a new task handler
func NewTaskHandler(taskUseCase *usecase.TaskUseCase) *TaskHandler {
	return &TaskHandler{
		taskUseCase: taskUseCase,
	}
}

// RegisterRoutes registers the task routes
func (h *TaskHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/tasks", h.handleTasks)
	mux.HandleFunc("/tasks/", h.handleTaskByID)
	mux.HandleFunc("/users/", h.handleTasksByUserID)
}

// handleTasks handles the /tasks endpoint
func (h *TaskHandler) handleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getTasks(w, r)
	case http.MethodPost:
		h.createTask(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleTaskByID handles the /tasks/{id} endpoint
func (h *TaskHandler) handleTaskByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	path := strings.TrimPrefix(r.URL.Path, "/tasks/")

	// Handle status update endpoints
	if strings.HasSuffix(path, "/in-progress") {
		id, err := strconv.ParseUint(strings.TrimSuffix(path, "/in-progress"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodPut {
			h.markTaskInProgress(w, r, id)
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if strings.HasSuffix(path, "/completed") {
		id, err := strconv.ParseUint(strings.TrimSuffix(path, "/completed"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodPut {
			h.markTaskCompleted(w, r, id)
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Handle regular task endpoints
	id, err := strconv.ParseUint(path, 10, 64)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getTaskByID(w, r, id)
	case http.MethodPut:
		h.updateTask(w, r, id)
	case http.MethodDelete:
		h.deleteTask(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleTasksByUserID handles the /users/{id}/tasks endpoint
func (h *TaskHandler) handleTasksByUserID(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from URL
	path := strings.TrimPrefix(r.URL.Path, "/users/")
	if !strings.Contains(path, "/tasks") {
		// Not a task endpoint, let the user handler handle it
		return
	}

	userIDStr := strings.Split(path, "/tasks")[0]
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodGet {
		h.getTasksByUserID(w, r, userID)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// getTasks handles GET /tasks
func (h *TaskHandler) getTasks(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10 // Default limit
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	offset := 0 // Default offset
	if offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// Get tasks
	tasks, err := h.taskUseCase.List(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, "Failed to get tasks: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return tasks
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tasks)
	if err != nil {
		return
	}
}

// getTasksByUserID handles GET /users/{id}/tasks
func (h *TaskHandler) getTasksByUserID(w http.ResponseWriter, r *http.Request, userID uint64) {
	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10 // Default limit
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	offset := 0 // Default offset
	if offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// Get tasks by user ID
	tasks, err := h.taskUseCase.GetByUserID(r.Context(), userID, limit, offset)
	if err != nil {
		http.Error(w, "Failed to get tasks: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Return tasks
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tasks)
	if err != nil {
		return
	}
}

// getTaskByID handles GET /tasks/{id}
func (h *TaskHandler) getTaskByID(w http.ResponseWriter, r *http.Request, id uint64) {
	// Get task
	task, err := h.taskUseCase.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	// Return task
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		return
	}
}

// createTask handles POST /tasks
func (h *TaskHandler) createTask(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req struct {
		Title       string     `json:"title"`
		Description string     `json:"description"`
		UserID      uint64     `json:"user_id"`
		DueDate     *time.Time `json:"due_date,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Title == "" || req.UserID == 0 {
		http.Error(w, "Title and user_id are required", http.StatusBadRequest)
		return
	}

	// Create task
	task, err := h.taskUseCase.Create(
		r.Context(),
		req.Title,
		req.Description,
		req.UserID,
		req.DueDate,
	)
	if err != nil {
		http.Error(w, "Failed to create task: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Return task
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		return
	}
}

// updateTask handles PUT /tasks/{id}
func (h *TaskHandler) updateTask(w http.ResponseWriter, r *http.Request, id uint64) {
	// Parse request body
	var req struct {
		Title       string            `json:"title"`
		Description string            `json:"description"`
		Status      entity.TaskStatus `json:"status"`
		DueDate     *time.Time        `json:"due_date,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	// Update task
	task, err := h.taskUseCase.Update(
		r.Context(),
		id,
		req.Title,
		req.Description,
		req.Status,
		req.DueDate,
	)
	if err != nil {
		http.Error(w, "Failed to update task: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Return task
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		return
	}
}

// deleteTask handles DELETE /tasks/{id}
func (h *TaskHandler) deleteTask(w http.ResponseWriter, r *http.Request, id uint64) {
	// Delete task
	if err := h.taskUseCase.Delete(r.Context(), id); err != nil {
		http.Error(w, "Failed to delete task: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Return success
	w.WriteHeader(http.StatusNoContent)
}

// markTaskInProgress handles PUT /tasks/{id}/in-progress
func (h *TaskHandler) markTaskInProgress(w http.ResponseWriter, r *http.Request, id uint64) {
	// Mark task as in progress
	task, err := h.taskUseCase.MarkInProgress(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to mark task as in progress: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Return task
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		return
	}
}

// markTaskCompleted handles PUT /tasks/{id}/completed
func (h *TaskHandler) markTaskCompleted(w http.ResponseWriter, r *http.Request, id uint64) {
	// Mark task as completed
	task, err := h.taskUseCase.MarkCompleted(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to mark task as completed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Return task
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		return
	}
}
