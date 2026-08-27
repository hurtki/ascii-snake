package lastaddr

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Repository struct {
	filePath string
}

func New(path string) (*Repository, error) {
	if path == "" {
		path = "~/ascii-snake/last-addr.txt"
	}

	resolvedPath, err := expandTilde(path)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path: %w", err)
	}

	return &Repository{filePath: resolvedPath}, nil
}

func (r *Repository) Load() (string, bool, error) {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", false, nil
		}
		return "", false, fmt.Errorf("read addr file: %w", err)
	}

	addr := strings.TrimSpace(string(data))
	if addr == "" {
		return "", false, nil
	}

	return addr, true, nil
}

func (r *Repository) Save(addr string) error {
	dir := filepath.Dir(r.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create dir %s: %w", dir, err)
	}

	err := os.WriteFile(r.filePath, []byte(strings.TrimSpace(addr)), 0644)
	if err != nil {
		return fmt.Errorf("write addr file: %w", err)
	}

	return nil
}

func expandTilde(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if path == "~" {
		return home, nil
	}

	if strings.HasPrefix(path, "~/") {
		return filepath.Join(home, path[2:]), nil
	}

	return path, nil
}
