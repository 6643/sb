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
	bodyV1 := []byte("package sb\n\nfunc demo() {}\n")
	if err := dir.SyncFingerprint("api.demo.go", bodyV1, 0644); err != nil {
		t.Fatalf("sync fingerprint v1 failed: %v", err)
	}

	contentV1, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fingerprinted file failed: %v", err)
	}
	if !bytes.Equal(contentV1, withFingerprint(bodyV1)) {
		t.Fatalf("expected initial fingerprinted content, got %q", string(contentV1))
	}

	bodyV2 := []byte("package sb\n\nfunc demo() { println(\"v2\") }\n")
	if err := dir.SyncFingerprint("api.demo.go", bodyV2, 0644); err != nil {
		t.Fatalf("sync fingerprint v2 failed: %v", err)
	}

	contentV2, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read updated fingerprinted file failed: %v", err)
	}
	if !bytes.Equal(contentV2, withFingerprint(bodyV2)) {
		t.Fatalf("expected managed file to update to new fingerprinted content, got %q", string(contentV2))
	}

	manualEdit := append([]byte{}, contentV2...)
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
