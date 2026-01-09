package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

var reservedNames = map[string]bool{
	"list":       true,
	"open":       true,
	"search":     true,
	"completion": true,
	"help":       true,
}

var validNameRegex = regexp.MustCompile(`^[a-zA-Z0-9-]+$`)

func GetStorageDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".remember")
}

func EnsureStorageDir() error {
	dir := GetStorageDir()
	if dir == "" {
		return fmt.Errorf("could not determine home directory")
	}
	return os.MkdirAll(dir, 0755)
}

func GetNotePath(name string) string {
	return filepath.Join(GetStorageDir(), name+".md")
}

func NoteExists(name string) bool {
	_, err := os.Stat(GetNotePath(name))
	return err == nil
}

func ValidateName(name string) error {
	if name == "" {
		return fmt.Errorf("note name cannot be empty")
	}
	if !validNameRegex.MatchString(name) {
		return fmt.Errorf("invalid name: use letters, numbers, hyphens only")
	}
	return nil
}

func IsReservedName(name string) bool {
	return reservedNames[name]
}

func ListNotes() ([]string, error) {
	dir := GetStorageDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var notes []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".md") {
			notes = append(notes, strings.TrimSuffix(name, ".md"))
		}
	}

	sort.Strings(notes)
	return notes, nil
}

func AppendToNote(name, content string) error {
	if err := EnsureStorageDir(); err != nil {
		return err
	}

	path := GetNotePath(name)
	timestamp := time.Now().Format("2006-01-02 15:04")

	var entry string
	if _, err := os.Stat(path); os.IsNotExist(err) {
		entry = fmt.Sprintf("## %s\n\n%s\n", timestamp, content)
	} else {
		entry = fmt.Sprintf("\n## %s\n\n%s\n", timestamp, content)
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(entry)
	return err
}
