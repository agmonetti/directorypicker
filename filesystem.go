package directorypicker

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// DirEntry represents a single filesystem entry in the picker listing.
type DirEntry struct {
	Name  string
	Path  string
	IsDir bool
}

func readDir(path string, showFiles bool) ([]DirEntry, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var result []DirEntry
	for _, e := range entries {
		name := e.Name()

		// Skip hidden entries
		if strings.HasPrefix(name, ".") {
			continue
		}

		// Resolve symlinks: check the target type
		isDir := e.IsDir()
		if e.Type()&os.ModeSymlink != 0 {
			info, err := os.Stat(filepath.Join(path, name))
			if err != nil {
				continue // broken symlink — skip
			}
			isDir = info.IsDir()
		}

		if isDir {
			result = append(result, DirEntry{
				Name:  name,
				Path:  filepath.Join(path, name),
				IsDir: true,
			})
		} else if showFiles {
			result = append(result, DirEntry{
				Name:  name,
				Path:  filepath.Join(path, name),
				IsDir: false,
			})
		}
	}

	// Sort: directories first, then files, alphabetically
	sort.Slice(result, func(i, j int) bool {
		a, b := result[i], result[j]
		if a.IsDir != b.IsDir {
			return a.IsDir
		}
		return strings.ToLower(a.Name) < strings.ToLower(b.Name)
	})

	return result, nil
}

func parentPath(path string) string {
	parent := filepath.Dir(path)
	if parent == path {
		return path
	}
	return parent
}

func homeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "/"
	}
	return home
}

func pathExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
