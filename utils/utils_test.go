package utils

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsFileNotExist(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T) string
		expected bool
	}{
		{
			"existing_file",
			func(t *testing.T) string {
				tempDir := t.TempDir()
				filePath := filepath.Join(tempDir, "existing.txt")
				err := os.WriteFile(filePath, []byte("test"), 0644)
				require.NoError(t, err)
				return filePath
			},
			false,
		},
		{
			"nonexistent_file",
			func(t *testing.T) string {
				tempDir := t.TempDir()
				return filepath.Join(tempDir, "nonexistent.txt")
			},
			true,
		},
		{
			"empty_path",
			func(t *testing.T) string {
				return ""
			},
			true,
		},
		{
			"invalid_path",
			func(t *testing.T) string {
				return "/nonexistent/directory/file.txt"
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := tt.setup(t)
			result := IsFileNotExist(filePath)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsDirectoryNotExist(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T) string
		expected bool
	}{
		{
			"existing_directory",
			func(t *testing.T) string {
				tempDir := t.TempDir()
				return tempDir
			},
			false,
		},
		{
			"nonexistent_directory",
			func(t *testing.T) string {
				tempDir := t.TempDir()
				return filepath.Join(tempDir, "nonexistent")
			},
			true,
		},
		{
			"empty_path",
			func(t *testing.T) string {
				return ""
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dirPath := tt.setup(t)
			result := IsDirectoryNotExist(dirPath)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestReplaceHomeDirWithTilde(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			"path_in_home",
			filepath.Join(homeDir, "Documents", "test.txt"),
			"~/Documents/test.txt",
		},
		{
			"exact_home_path",
			homeDir,
			"~",
		},
		{
			"path_outside_home",
			"/usr/local/bin/test",
			"/usr/local/bin/test",
		},
		{
			"empty_path",
			"",
			"",
		},
		{
			"relative_path",
			"./relative/path",
			"./relative/path",
		},
		{
			"path_with_home_substring",
			"/some" + homeDir + "/fake",
			"/some" + homeDir + "/fake",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ReplaceHomeDirWithTilde(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPrintTable(t *testing.T) {
	var buf bytes.Buffer
	headers := []string{"Name", "Type", "Status"}
	rows := [][]string{
		{"key1", "ed25519", "active"},
		{"key2", "rsa", "inactive"},
		{"key3", "ed25519", "active"},
	}

	PrintTable(&buf, headers, rows)

	output := buf.String()

	// headers
	assert.Contains(t, output, "Name")
	assert.Contains(t, output, "Type")
	assert.Contains(t, output, "Status")

	// data rows
	assert.Contains(t, output, "key1")
	assert.Contains(t, output, "key2")
	assert.Contains(t, output, "key3")
	assert.Contains(t, output, "ed25519")
	assert.Contains(t, output, "rsa")
	assert.Contains(t, output, "active")
	assert.Contains(t, output, "inactive")

	lines := strings.Split(output, "\n")
	assert.Greater(t, len(lines), 3)
}

func TestPrintTable_EmptyData(t *testing.T) {
	var buf bytes.Buffer
	headers := []string{"Name", "Type", "Status"}
	rows := [][]string{}

	PrintTable(&buf, headers, rows)

	output := buf.String()

	// headers
	assert.Contains(t, output, "Name")
	assert.Contains(t, output, "Type")
	assert.Contains(t, output, "Status")

	lines := strings.Split(strings.TrimSpace(output), "\n")
	assert.Equal(t, 1, len(lines))
}
