package loader

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTemp: %v", err)
	}
	return p
}

func TestLoadEnvFile_Basic(t *testing.T) {
	path := writeTemp(t, "APP_ENV=production\nDB_HOST=localhost\n")
	got, err := LoadEnvFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["APP_ENV"] != "production" {
		t.Errorf("APP_ENV: want production, got %q", got["APP_ENV"])
	}
	if got["DB_HOST"] != "localhost" {
		t.Errorf("DB_HOST: want localhost, got %q", got["DB_HOST"])
	}
}

func TestLoadEnvFile_CommentsAndBlanks(t *testing.T) {
	content := "# comment\n\nFOO=bar\n"
	path := writeTemp(t, content)
	got, err := LoadEnvFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got["FOO"] != "bar" {
		t.Errorf("unexpected map: %v", got)
	}
}

func TestLoadEnvFile_InvalidLine(t *testing.T) {
	path := writeTemp(t, "INVALID_LINE\n")
	_, err := LoadEnvFile(path)
	if err == nil {
		t.Fatal("expected error for invalid line, got nil")
	}
}

func TestLoadEnvFile_NotFound(t *testing.T) {
	_, err := LoadEnvFile("/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadEnvFile_ValueWithEquals(t *testing.T) {
	path := writeTemp(t, "JDBC_URL=jdbc:mysql://host/db?user=root&pass=s=ecret\n")
	got, err := LoadEnvFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "jdbc:mysql://host/db?user=root&pass=s=ecret"
	if got["JDBC_URL"] != want {
		t.Errorf("JDBC_URL: want %q, got %q", want, got["JDBC_URL"])
	}
}
