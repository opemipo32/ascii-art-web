package banner_test

import (
	"ascii-art/banner"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempBanner(t *testing.T, name, content string) string {
	t.Helper()
	dir := t.TempDir()

	bannersDir := filepath.Join(dir, "banner")
	if err := os.MkdirAll(bannersDir, 0o755); err != nil {
		t.Fatalf("failed to create banner dir: %v", err)
	}

	path := filepath.Join(bannersDir, name+".txt")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write temp banner: %v", err)
	}

	orig, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	return path
}

func buildMinimalBanner() string {
	var sb strings.Builder
	for i := 0; i < 95; i++ {
		sb.WriteString("\n")
		for row := 0; row < 8; row++ {
			sb.WriteString("----\n")
		}
	}
	return sb.String()
}

func TestLoad_Valid(t *testing.T) {
	content := buildMinimalBanner()
	writeTempBanner(t, "test", content)

	chars, err := banner.Load("test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(chars) != 95 {
		t.Errorf("expected 95 chars, got %d", len(chars))
	}
	for i, ch := range chars {
		if len(ch) != 8 {
			t.Errorf("char %d: expected 8 rows, got %d", i, len(ch))
		}
	}
}

func TestLoad_MissingFile(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	_ = os.Chdir(dir)
	t.Cleanup(func() { _ = os.Chdir(orig) })

	_ = os.MkdirAll(filepath.Join(dir, "banners"), 0o755)

	_, err := banner.Load("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing banner, got nil")
	}
}

func TestLoad_MalformedFile(t *testing.T) {
	var sb strings.Builder
	for i := 0; i < 10; i++ {
		sb.WriteString("\n")
		for row := 0; row < 8; row++ {
			sb.WriteString("----\n")
		}
	}
	writeTempBanner(t, "bad", sb.String())

	_, err := banner.Load("bad")
	if err == nil {
		t.Fatal("expected error for malformed banner, got nil")
	}
}
