package repository

import (
	"database/sql"
	"fmt"

	"github.com/ybotet/pz6_csrf_xss/shared/models"
)

type TaskRepository interface {
    GetAll() ([]models.Task, error)
    GetByID(id string) (*models.Task, error)
    Create(task *models.Task) error
    Update(task *models.Task) error
    Delete(id string) error
    SearchByTitle(title string) ([]models.Task, error)  // Endpoint vulnerable
}

type PostgresTaskRepository struct {
    db *sql.DB
}

func NewPostgresTaskRepository(db *sql.DB) *PostgresTaskRepository {
    return &PostgresTaskRepository{db: db}
}

// VERSIÓN VULNERABLE (para demostración)
func (r *PostgresTaskRepository) SearchByTitleVulnerable(title string) ([]models.Task, error) {
    //  MAL: Concatenación de strings - VULNERABLE A SQL INJECTION
    query := fmt.Sprintf("SELECT id, title, description, done, created_at FROM tasks WHERE title = '%s'", title)
    
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var tasks []models.Task
    for rows.Next() {
        var t models.Task
        err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Done, &t.CreatedAt)
        if err != nil {
            return nil, err
        }
        tasks = append(tasks, t)
    }
    return tasks, nil
}

// VERSIÓN SEGURA (parametrizada)
func (r *PostgresTaskRepository) SearchByTitle(title string) ([]models.Task, error) {
    //  BIEN: Consulta parametrizada - SEGURA
    query := "SELECT id, title, description, done, created_at FROM tasks WHERE title = "
    
    rows, err := r.db.Query(query, title)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var tasks []models.Task
    for rows.Next() {
        var t models.Task
        err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Done, &t.CreatedAt)
        if err != nil {
            return nil, err
        }
        tasks = append(tasks, t)
    }
    return tasks, nil
}

// Las demás implementaciones...
func (r *PostgresTaskRepository) GetAll() ([]models.Task, error) {
    rows, err := r.db.Query("SELECT id, title, description, done, created_at FROM tasks ORDER BY created_at DESC")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var tasks []models.Task
    for rows.Next() {
        var t models.Task
        err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Done, &t.CreatedAt)
        if err != nil {
            return nil, err
        }
        tasks = append(tasks, t)
    }
    return tasks, nil
}

func (r *PostgresTaskRepository) GetByID(id string) (*models.Task, error) {
    var t models.Task
    err := r.db.QueryRow("SELECT id, title, description, done, created_at FROM tasks WHERE id = ", id).
        Scan(&t.ID, &t.Title, &t.Description, &t.Done, &t.CreatedAt)
    if err != nil {
        return nil, err
    }
    return &t, nil
}

func (r *PostgresTaskRepository) Create(task *models.Task) error {
    query := `INSERT INTO tasks (id, title, description, done, created_at) 
              VALUES ($1, $2, $3, $4, $5)`
    
    _, err := r.db.Exec(query, 
        task.ID, 
        task.Title, 
        task.Description, 
        task.Done, 
        task.CreatedAt)
    
    if err != nil {
        return fmt.Errorf("error inserting task: %w", err)
    }
    return err
}

func (r *PostgresTaskRepository) Update(task *models.Task) error {
    query := "UPDATE tasks SET title=, description=, done= WHERE id="
    _, err := r.db.Exec(query, task.Title, task.Description, task.Done, task.ID)
    return err
}

func (r *PostgresTaskRepository) Delete(id string) error {
    _, err := r.db.Exec("DELETE FROM tasks WHERE id = ", id)
    return err
}
