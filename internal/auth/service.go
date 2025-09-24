package auth 

import (
	"fmt"
)

type Service struct {
	storage *Storage
}

func NewService(storage *Storage) *Service {
	return &Service{storage: storage}
}

// RegisterAPIKey ganerates a new key and saves it hashed in db
func (s *Service) RegisterAPIKey(owner string) (string, error) {
	plain, hash, err := GenerateAPIKey()
	if err != nil {
		return "", fmt.Errorf("failed to generate api key: %w", err)
	}


	_, err = s.storage.AddAPIKey(hash, owner)
	if err != nil {
		return "", fmt.Errorf("failed to save key: %w", err)
	}
	return plain, nil
}

// CheckPermission verifies if apiKey has given permission
func (s *Service) CheckPermission(apiKeyPlain string, permission string) (bool, error) {
	hash := HashKey(apiKeyPlain)
	return s.storage.HasPermission(hash, permission)
}

// ListKeys return all stored keys
func (s *Service) ListKeys() ([]APIKey, error) {
	return s.storage.ListAPIKeys()
}

// CreatePermission creates a new permission
func (s *Service) CreatePermission(name string) error {
	_, err := s.storage.db.Exec("INSERT INTO permissions(name) VALUES(?)", name)
	if err != nil {
		return fmt.Errorf("failed to create permission: %w", err)
	}
	return nil
}

// GrantPermission assings a permission to a key
func (s *Service) GrantPermission(keyID, permID int64) error {
	_, err := s.storage.db.Exec(
		"INSERT OR IGNORE INTO api_key_permissions(api_key_id, permission_id) VALUES (?, ?)",
		keyID, permID,
	)
	if err != nil {
		return fmt.Errorf("failed to grant permission: %w", err)
	}
	return nil
}

func (s *Service) ListPermissions() ([]Permission, error) {
    rows, err := s.storage.db.Query("SELECT id, name FROM permissions")
    if err != nil {
        return nil, fmt.Errorf("failed to list permissions: %w", err)
    }
    defer rows.Close()

    var perms []Permission
    for rows.Next() {
        var p Permission
        if err := rows.Scan(&p.ID, &p.Name); err != nil {
            return nil, err
        }
        perms = append(perms, p)
    }
    return perms, nil
}