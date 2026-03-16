package handlers

import (
    "encoding/json"
    "net/http"
    "time"
    
    "github.com/google/uuid"
    "github.com/ybotet/pz6_csrf_xss/services/tasks/internal/repository"
    "github.com/ybotet/pz6_csrf_xss/shared/models"
)

type TaskHandler struct {
    repo repository.TaskRepository
}

func NewTaskHandler(repo repository.TaskRepository) *TaskHandler {
    return &TaskHandler{repo: repo}
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
    tasks, err := h.repo.GetAll()
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
    var task models.Task
    if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    task.ID = uuid.New().String()
    task.CreatedAt = time.Now()
    task.Done = false
    
    if err := h.repo.Create(&task); err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(task)
}

// ENDPOINT VULNERABLE (para demostración)
func (h *TaskHandler) SearchTasksVulnerable(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Query().Get("title")
    if title == "" {
        http.Error(w, "title parameter required", http.StatusBadRequest)
        return
    }
    
    // Usar versión vulnerable (solo para demo)
    repo, ok := h.repo.(*repository.PostgresTaskRepository)
    if !ok {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    
    tasks, err := repo.SearchByTitleVulnerable(title)
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    
    json.NewEncoder(w).Encode(tasks)
}

// ENDPOINT SEGURO (parametrizado)
func (h *TaskHandler) SearchTasks(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Query().Get("title")
    if title == "" {
        http.Error(w, "title parameter required", http.StatusBadRequest)
        return
    }
    
    tasks, err := h.repo.SearchByTitle(title)
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    
    json.NewEncoder(w).Encode(tasks)
}
