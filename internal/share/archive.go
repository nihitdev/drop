package share

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func zipDirectory(source string) (archivePath string, retErr error) {
	tmp, err := os.CreateTemp("", "drop-*.zip")
	if err != nil {
		return "", fmt.Errorf("create temporary file: %w", err)
	}
	archivePath = tmp.Name()
	zipWriter := zip.NewWriter(tmp)

	defer func() {
		if retErr != nil {
			_ = zipWriter.Close()
			_ = tmp.Close()
			_ = os.Remove(archivePath)
		}
	}()

	err = filepath.Walk(source, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() {
			return nil
		}
		if !info.Mode().IsRegular() {
			return nil
		}

		rel, err := filepath.Rel(source, path)
		if err != nil {
			return fmt.Errorf("make archive path for %q: %w", path, err)
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("create ZIP header for %q: %w", path, err)
		}
		header.Name = filepath.ToSlash(rel)
		header.Method = zip.Deflate
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("add %q to archive: %w", path, err)
		}

		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("open %q: %w", path, err)
		}
		_, copyErr := io.Copy(writer, file)
		closeErr := file.Close()
		if copyErr != nil {
			return fmt.Errorf("archive %q: %w", path, copyErr)
		}
		if closeErr != nil {
			return fmt.Errorf("close %q: %w", path, closeErr)
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if err := zipWriter.Close(); err != nil {
		return "", fmt.Errorf("finish ZIP archive: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return "", fmt.Errorf("close temporary archive: %w", err)
	}
	return archivePath, nil
}
