package share

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type Target struct {
	Path        string
	DisplayName string
	Size        int64
	Temporary   bool
	cleanupOnce sync.Once
}

func Prepare(path string) (*Target, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("cannot access %q: %w", path, err)
	}

	displayName := info.Name()
	preparedPath := path
	temporary := false
	if info.IsDir() {
		preparedPath, err = zipDirectory(path)
		if err != nil {
			return nil, fmt.Errorf("failed to create temporary archive: %w", err)
		}
		displayName = filepath.Base(filepath.Clean(path)) + ".zip"
		temporary = true
	} else if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("%q is not a regular file or directory", path)
	}

	absPath, err := filepath.Abs(preparedPath)
	if err != nil {
		if temporary {
			_ = os.Remove(preparedPath)
		}
		return nil, fmt.Errorf("resolve absolute path: %w", err)
	}
	preparedInfo, err := os.Stat(absPath)
	if err != nil {
		if temporary {
			_ = os.Remove(preparedPath)
		}
		return nil, fmt.Errorf("inspect prepared target: %w", err)
	}

	return &Target{Path: absPath, DisplayName: displayName, Size: preparedInfo.Size(), Temporary: temporary}, nil
}

func (t *Target) Cleanup() {
	if t == nil || !t.Temporary {
		return
	}
	t.cleanupOnce.Do(func() { _ = os.Remove(t.Path) })
}
