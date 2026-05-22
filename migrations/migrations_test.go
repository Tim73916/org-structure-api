package migrations

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigrationFilesExist(t *testing.T) {
	files, err := filepath.Glob("*.sql")
	assert.NoError(t, err)
	assert.NotEmpty(t, files, "No migration files found")
}

func TestMigrationFilesHaveUpAndDown(t *testing.T) {
	files, _ := filepath.Glob("*.sql")

	for _, file := range files {
		content, err := os.ReadFile(file)
		assert.NoError(t, err)

		hasUp := strings.Contains(string(content), "-- +goose Up")
		hasDown := strings.Contains(string(content), "-- +goose Down")

		assert.True(t, hasUp, "File %s missing +goose Up", file)
		assert.True(t, hasDown, "File %s missing +goose Down", file)
	}
}
