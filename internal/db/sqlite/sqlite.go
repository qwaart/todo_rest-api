package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

type Storage struct {
	db *sql.DB
}

// New creates a new SQLite storage and initializes the todo table.
func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Create table for tasks only
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS todo(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task TEXT NOT NULL,
		completed BOOLEAN DEFAULT FALSE
	);

	CREATE INDEX IF NOT EXISTS idx_task ON todo(task);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// DB exposes raw *sql.DB for other packages (e.g. auth)
func (s *Storage) DB() *sql.DB {
	return s.db
}

// AddTask adds a new task and returns its ID
func (s *Storage) AddTask(task string) (int64, error) {
	const op = "storage.sqlite.AddTask"

	stmt, err := s.db.Prepare("INSERT INTO todo(task) VALUES(?)")
	if err != nil {
		return 0, fmt.Errorf("%s: Prepare: %w", op, err)
	}
	defer stmt.Close()

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

// GetTaskByID returns task text and completed status by ID
func (s *Storage) GetTaskByID(id int64) (string, bool, error) {
	const op = "storage.sqlite.GetTaskByID"

	var task string
	var completed bool

	err := s.db.QueryRow("SELECT task, completed FROM todo WHERE id = ?", id).
		Scan(&task, &completed)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", false, fmt.Errorf("No task with id=%d: %w", id, sql.ErrNoRows)
		}
		return "", false, fmt.Errorf("%s: QueryRow: %w", op, err)
	}

	return task, completed, nil
}

func (s *Storage) DeleteTaskByID(id int64) error {
	_, err := s.db.Exec(`DELETE FROM todo WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	return nil
}

func (s *Storage) MarkTaskTrue(id int64) error {
	_, err := s.db.Exec(`UPDATE todo SET completed = 1 WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to mark task completed: %w", err)
	}
	return nil
}

func (s *Storage) MarkTaskFalse(id int64) error {
	_, err := s.db.Exec(`UPDATE todo SET completed = 0 WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to mark task uncompleted: %w", err)
	}
	return nil
}
