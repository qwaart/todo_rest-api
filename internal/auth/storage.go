package auth

import (
	"database/sql"
	"fmt"
	"log/slog"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

// Setup tables
func (s *Storage) Init() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS api_keys(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			key_hash TEXT NOT NULL UNIQUE,
			owner TEXT NOT NULL,
			revoked BOOLEAN DEFAULT 0
		);

		CREATE TABLE IF NOT EXISTS permissions(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE
		);

		CREATE TABLE IF NOT EXISTS api_key_permissions(
			api_key_id INTEGER,
			permission_id INTEGER,
			PRIMARY KEY (api_key_id, permission_id),
			FOREIGN KEY (api_key_id) REFERENCES api_keys(id),
			FOREIGN KEY (permission_id) REFERENCES permissions(id)
		);
		`)
		if err != nil {
			return fmt.Errorf("failed to init auth tables: %w", err)
		}
		return nil
}

func (s *Storage) ListAPIKeys() ([]APIKey, error) {
	rows, err := s.db.Query("SELECT id, key_hash, owner, revoked FROM api_keys")
	if err != nil {
		return nil, fmt.Errorf("failed to list keys: %w", err)
	}
	defer rows.Close()

	var keys []APIKey
	for rows.Next() {
		var k APIKey
		if err := rows.Scan(&k.ID, &k.Key, &k.Owner, &k.Revoked); err != nil {
			return nil, err
		}
		keys = append(keys, k)
	}
	return keys, nil
}

// permissions
func (s *Storage) HasPermission(providedKey, permission string) (bool, error) {
	hashed := HashKey(providedKey)

	const query = `
	SELECT COUNT(*)
	FROM api_keys k
	JOIN api_key_permissions kp ON k.id = kp.api_key_id
	JOIN permissions p ON kp.permission_id = p.id
	WHERE k.key_hash = ? AND p.name = ? AND k.revoked = 0
	`

	var count int
	err := s.db.QueryRow(query, hashed, permission).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check permission: %w", err)
	}

	return count > 0, nil
}

// EnsureAdminSetup checks for the presence of the admin key and rights
// and creates them only if they are missing.
func (s *Storage) EnsureAdminSetup(log *slog.Logger) error {
	const op = "auth.storage.EnsureAdminSetup"

	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM api_keys WHERE owner = 'admin'").Scan(&count)
	if err != nil {
		return fmt.Errorf("%s: failed to check admin key: %w", op, err)
	}

	if count > 0 {
		log.Info("Admin key already exists")
		return nil
	}

	plain, hash, err := GenerateAPIKey()
	if err != nil {
		return fmt.Errorf("%s: failed to generate admin key: %w", op, err)
	}

	res, err := s.db.Exec("INSERT INTO api_keys(key_hash, owner) VALUES (?, 'admin')", hash)
	if err != nil {
		return fmt.Errorf("%s: failed to insert admin key: %w", op, err)
	}

	// create a permission admin
	if err := s.CreatePermission("admin"); err != nil {
		return fmt.Errorf("%s: failed to create admin permission: %w", op, err)
	}

	keyID, _ := res.LastInsertId()
	if err := s.GrantPermission(fmt.Sprint(keyID), "admin"); err != nil {
		return fmt.Errorf("%s: failed to grant admin permission: %w", op, err)
	}

	log.Warn("Admin API key generated - save it securely, it will not be sown again ",
		slog.String("key", plain))

	return nil
}

//RegisterKey create new api-key and save it to db
func(s *Storage) RegisterKey(owner string) (string, error) {
	plainKey, hashedKey, err := GenerateAPIKey()
	if err != nil {
		return "", fmt.Errorf("failed to generate key: %w", err)
	}

	_, err = s.db.Exec("INSERT INTO api_keys (key_hash, owner) VALUES (?, ?)", hashedKey, owner)
	if err != nil {
		return "", fmt.Errorf("failed to save api key: %w", err)
	}

	return plainKey, nil //return only the plainkey to the user
}

// ValidateKey checks if the key is in the database and has not been revoked
func (s *Storage) ValidateKey(providedKey string) (bool, error) {
	hashed := HashKey(providedKey)

	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM api_keys WHERE key_hash = ? AND revoked = 0)", hashed).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check api key: %w", err)
	}

	return exists, nil
}

// GrantPermission adds a right to the API key
func (s *Storage) GrantPermission(apiKeyID string, permissionName string) error {
	var permID int
	err := s.db.QueryRow("SELECT id FROM permissions WHERE name = ?", permissionName).Scan(&permID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("permission not found")
	} else if err != nil {
		return fmt.Errorf("failed to lookup permission: %w", err)
	}

	_, err = s.db.Exec(`
		INSERT OR IGNORE INTO api_key_permissions(api_key_id, permission_id)
		VALUES (?, ?)`, apiKeyID, permID)
	if err != nil {
		return fmt.Errorf("failed to grant permission: %w", err)
	}
	return nil
}

// CreatePermission create a new permission, if they not exist
func (s *Storage) CreatePermission(name string) error {
	_, err := s.db.Exec("INSERT OR IGNORE INTO permissions(name) VALUES (?)", name)
	if err != nil {
		return fmt.Errorf("failed to create permission: %w", err)
	}
	return nil
}

// AddAPIKey inserts hashedKey into api_keys and returns inserted id
func (s *Storage) AddAPIKey(hashedKey, owner string) (int64, error) {
	const op = "auth.storage.AddAPIKey"

	res, err := s.db.Exec("INSERT INTO api_keys(key_hash, owner) VALUES (?, ?)", hashedKey, owner)
	if err != nil {
		return 0, fmt.Errorf("%s: insert failed: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: get lastInsertId: %w", op, err)
	}
	return id, nil
}
