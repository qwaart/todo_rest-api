package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Storage is a wrapper over sql.DB for workin with todo task
type Storage struct {
	db *sql.DB
}

// New creates a new SQLite storage and initializes the todo table.
// If the database file does not exist, SQLite will create it automatically.
func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Initializing the table and index
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS todo(
		id INTEGER PRIMARY KEY,
		task TEXT NOT NULL,
		completed BOOLEAN DEFAULT FALSE);
	CREATE INDEX IF NOT EXISTS idx_task ON todo(task);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// POST - AddTask adds a new task to the db and return its id
func (s *Storage) AddTask(task string) (int64, error) {
	const op = "storage.sqlite.AddTask"

	stmt, err := s.db.Prepare("INSERT INTO todo(task) VALUES(?)")
	if err != nil {
		return 0, fmt.Errorf("%s: Prepare: %w", op, err)
	}
	defer stmt.Close() // close the prepared statement

	result, err := stmt.Exec(task)
	if err != nil {
		return 0, fmt.Errorf("%s: Exec: %w", op, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: lastInsertId: %w", op, err)
	}

	return id, nil
}

 // GET - return a task by id
 // if task not found, it return an error
func (s *Storage) GetTaskByID(id int64) (string, bool, error) {
	const op = "storage.sqlite.GetTaskByID"
	var task string
	var completed bool

	err := s.db.QueryRow("SELECT task, completed FROM todo WHERE id = ?", id).Scan(&task, &completed)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", false, fmt.Errorf("No task with id=%d: %w", id, sql.ErrNoRows)
		}
		return "", false, fmt.Errorf("%s: QueryRow: %w", op, err)
	}

	return task, completed, nil
}

 // DELETE - delete task by id :)
func (s *Storage) DeleteTaskByID(id int64) error {
	query := `DELETE FROM todo WHERE id = ?`

	_, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

 // PATCH - marks the task as completed/uncompleted
func (s *Storage) MarkTaskTrue(id int64) error {
	query := `UPDATE todo SET completed = 1 WHERE id = ?`

	_, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("an error when marking a task as completed: %w", err)
	}

	return nil
}

func (s *Storage) MarkTaskFalse(id int64) error {
	query := `UPDATE todo SET completed = 0 WHERE id = ?`

	_, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("an error when marking a task as uncompleted: %w", err)
	}

	return nil
}
