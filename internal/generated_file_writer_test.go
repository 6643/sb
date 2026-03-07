package internal

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGeneratedDirWriteCreatesDirsAndSkipsUnchangedContent(t *testing.T) {
	dir, err := newGeneratedDir(filepath.Join(t.TempDir(), "out"))
	if err != nil {
		t.Fatalf("create generated dir failed: %v", err)
	}

	path := filepath.Join(dir.root, "nested", "file.txt")
	if err := dir.Write("nested/file.txt", []byte("alpha"), 0644); err != nil {
		t.Fatalf("write generated file failed: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read generated file failed: %v", err)
	}
	if string(content) != "alpha" {
		t.Fatalf("expected generated content alpha, got %q", string(content))
	}

	stableTime := time.Date(2001, time.February, 3, 4, 5, 6, 0, time.UTC)
	if err := os.Chtimes(path, stableTime, stableTime); err != nil {
		t.Fatalf("set generated file time failed: %v", err)
	}
	if err := dir.Write("nested/file.txt", []byte("alpha"), 0644); err != nil {
		t.Fatalf("rewrite unchanged generated file failed: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat generated file failed: %v", err)
	}
	if !info.ModTime().Equal(stableTime) {
		t.Fatalf("expected unchanged file mod time %s, got %s", stableTime, info.ModTime())
	}

	if err := dir.WriteIfAbsent("nested/file.txt", []byte("beta"), 0644); err != nil {
		t.Fatalf("write-if-absent existing file failed: %v", err)
	}
	content, err = os.ReadFile(path)
	if err != nil {
		t.Fatalf("read generated file after write-if-absent failed: %v", err)
	}
	if string(content) != "alpha" {
		t.Fatalf("expected existing generated file to stay alpha, got %q", string(content))
	}
}

func TestGeneratedDirSyncFingerprintOnlyUpdatesManagedFiles(t *testing.T) {
	dir, err := newGeneratedDir(filepath.Join(t.TempDir(), "out"))
	if err != nil {
		t.Fatalf("create generated dir failed: %v", err)
	}

	path := filepath.Join(dir.root, "api.demo.go")
	initialBody := []byte("package sb\n\nfunc demo() {}\n")
	if err := dir.SyncFingerprint("api.demo.go", initialBody, 0644); err != nil {
		t.Fatalf("sync fingerprint initial failed: %v", err)
	}

	initialContent, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fingerprinted file failed: %v", err)
	}
	if !bytes.Equal(initialContent, withFingerprint(initialBody)) {
		t.Fatalf("expected initial fingerprinted content, got %q", string(initialContent))
	}

	updatedBody := []byte("package sb\n\nfunc demo() { println(\"updated\") }\n")
	if err := dir.SyncFingerprint("api.demo.go", updatedBody, 0644); err != nil {
		t.Fatalf("sync fingerprint update failed: %v", err)
	}

	updatedContent, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read updated fingerprinted file failed: %v", err)
	}
	if !bytes.Equal(updatedContent, withFingerprint(updatedBody)) {
		t.Fatalf("expected managed file to update to new fingerprinted content, got %q", string(updatedContent))
	}

	manualEdit := append([]byte{}, updatedContent...)
	manualEdit = append(manualEdit, []byte("// manual edit\n")...)
	if err := os.WriteFile(path, manualEdit, 0644); err != nil {
		t.Fatalf("write manual edit failed: %v", err)
	}

	if err := dir.SyncFingerprint("api.demo.go", []byte("package sb\n\nfunc demo() { println(\"v3\") }\n"), 0644); err != nil {
		t.Fatalf("sync fingerprint after manual edit failed: %v", err)
	}

	finalContent, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fingerprinted file after manual edit failed: %v", err)
	}
	if !bytes.Equal(finalContent, manualEdit) {
		t.Fatalf("expected manual edit to be preserved, got %q", string(finalContent))
	}
}
