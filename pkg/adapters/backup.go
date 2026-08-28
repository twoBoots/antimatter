package adapters

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// DefaultBackupDir returns the default backup base path (~/.antimatter/backups)
func DefaultBackupDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", ".antimatter", "backups")
	}
	return filepath.Join(home, ".antimatter", "backups")
}

// CreateBackup copies the source configuration file to a timestamped backup destination
func CreateBackup(sourcePath string, agentName string) (string, error) {
	return CreateBackupWithBaseDir(sourcePath, agentName, DefaultBackupDir())
}

// CreateBackupWithBaseDir copies the source file to a custom backup base directory
func CreateBackupWithBaseDir(sourcePath string, agentName string, baseDir string) (string, error) {
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return "", nil // No existing file to back up
	}

	targetDir := filepath.Join(baseDir, agentName)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create backup directory '%s': %w", targetDir, err)
	}

	timestamp := time.Now().UTC().Format("20060102_150405")
	baseName := filepath.Base(sourcePath)
	destPath := filepath.Join(targetDir, fmt.Sprintf("%s_%s", timestamp, baseName))

	srcFile, err := os.Open(sourcePath)
	if err != nil {
		return "", fmt.Errorf("failed to open source config '%s': %w", sourcePath, err)
	}
	defer srcFile.Close()

	destFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to create backup file '%s': %w", destPath, err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		return "", fmt.Errorf("failed to copy content to backup '%s': %w", destPath, err)
	}

	return destPath, nil
}

// ListBackups returns all created backup paths for a given agent, sorted newest first
func ListBackups(agentName string) ([]string, error) {
	return ListBackupsWithBaseDir(agentName, DefaultBackupDir())
}

// ListBackupsWithBaseDir lists backups from a specified base directory
func ListBackupsWithBaseDir(agentName string, baseDir string) ([]string, error) {
	targetDir := filepath.Join(baseDir, agentName)
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	entries, err := os.ReadDir(targetDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup directory '%s': %w", targetDir, err)
	}

	var results []string
	for _, entry := range entries {
		if !entry.IsDir() {
			results = append(results, filepath.Join(targetDir, entry.Name()))
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i] > results[j] // descending order (newest first)
	})

	return results, nil
}
