package share

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

func TestPrepareDirectory(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "nested"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "nested", "hello.txt"), []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	target, err := Prepare(dir)
	if err != nil {
		t.Fatal(err)
	}
	archivePath := target.Path
	defer target.Cleanup()
	if !target.Temporary || target.DisplayName != filepath.Base(dir)+".zip" {
		t.Fatalf("unexpected target: %+v", target)
	}
	zr, err := zip.OpenReader(target.Path)
	if err != nil {
		t.Fatal(err)
	}
	defer zr.Close()
	if len(zr.File) != 1 || zr.File[0].Name != "nested/hello.txt" {
		t.Fatalf("unexpected ZIP entries: %+v", zr.File)
	}
	target.Cleanup()
	if _, err := os.Stat(archivePath); !os.IsNotExist(err) {
		t.Fatalf("temporary archive was not removed: %v", err)
	}
}
