package sqlite

import (
	"database/sql"
	"fmt"

	_ "rest_api/internal/db"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

// This func create storage
func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS todo(
		id INTEGER PRIMARY KEY,
		task TEXT NOT NULL,
		completed BOOLEAN DEFAULT FALSE,
		private BOOLEAN DEFAULT FALSE);
	CREATE INDEX IF NOT EXISTS idx_task ON todo(task);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// POST
func (s *Storage) AddTask(task string, private bool) (int64, error) {
	const op = "storage.sqlite.AddTask"

	stmt, err := s.db.Prepare("INSERT INTO todo(task, private) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: Prepare: %w", op, err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(task, private)
	if err != nil {
		return 0, fmt.Errorf("%s: Exec: %w", op, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: lastInsertId: %w", op, err)
	}

	return id, nil
}
 // GET
func (s *Storage) GetTaskByID(id int64) (string, bool, bool, error) {
	const op = "storage.sqlite.GetTaskByID"
	var task string
	var completed bool
	var private bool

	err := s.db.QueryRow("SELECT task, completed, private FROM todo WHERE id = ?", id).Scan(&task, &completed, &private)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", false, false, fmt.Errorf("No task with id=%d", id)
		}
		return "", false, false, fmt.Errorf("%s: QuerryRow: %w", op, err)
	}

	return task, completed, private, nil
}

 // DELETE
func (s *Storage) DeleteTaskByID(id int64) error {
	query := `DELETE FROM todo WHERE id = ?`

	_, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

 // UPDATE/PATCH
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


 //need for other handlers
func (s *Storage) GetTaskFullByID(id int64) (string, bool, bool, error) {
	const op = "storage.sqlite.GetTaskFullByID"
	var task string
	var completed bool
	var private bool

 	err := s.db.QueryRow("SELECT task, completed, private FROM todo where id = ?", id).Scan(&task, &completed, &private)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", false, false, fmt.Errorf("no task with id=%d", id)
		}
		return "", false, false, fmt.Errorf("%s: QuerryRow: %w", op, err)
	}

	return task, completed, private, nil
}