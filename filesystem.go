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
	Size  int64 // used for SortBySizeDesc; 0 for directories
}

func readDir(path string, showFiles bool, showHidden bool, allowedExts []string, sortOrder SortOrder) ([]DirEntry, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var result []DirEntry
	for _, e := range entries {
		name := e.Name()

		// Skip hidden entries if not enabled
		if !showHidden && strings.HasPrefix(name, ".") {
			continue
		}

		// Resolve symlinks: check the target type
		isDir := e.IsDir()
		var fileSize int64
		if e.Type()&os.ModeSymlink != 0 {
			info, err := os.Stat(filepath.Join(path, name))
			if err != nil {
				continue // broken symlink — skip
			}
			isDir = info.IsDir()
			if !isDir {
				fileSize = info.Size()
			}
		} else if !isDir {
			if info, err := e.Info(); err == nil {
				fileSize = info.Size()
			}
		}

		if isDir {
			result = append(result, DirEntry{
				Name:  name,
				Path:  filepath.Join(path, name),
				IsDir: true,
			})
		} else if showFiles {
			// Apply extension filter when a non-empty, non-wildcard list is provided
			if len(allowedExts) > 0 && allowedExts[0] != "*" {
				ext := strings.ToLower(filepath.Ext(name))
				matched := false
				for _, allowed := range allowedExts {
					if ext == strings.ToLower(allowed) {
						matched = true
						break
					}
				}
				if !matched {
					continue
				}
			}
			result = append(result, DirEntry{
				Name:  name,
				Path:  filepath.Join(path, name),
				IsDir: false,
				Size:  fileSize,
			})
		}
	}

	// Sort according to SortOrder
	sort.SliceStable(result, func(i, j int) bool {
		a, b := result[i], result[j]
		switch sortOrder {
		case SortByNameAsc:
			return strings.ToLower(a.Name) < strings.ToLower(b.Name)
		case SortByNameDesc:
			return strings.ToLower(a.Name) > strings.ToLower(b.Name)
		case SortFilesFirst:
			if a.IsDir != b.IsDir {
				return !a.IsDir // files first
			}
			return strings.ToLower(a.Name) < strings.ToLower(b.Name)
		case SortBySizeDesc:
			if a.IsDir != b.IsDir {
				return a.IsDir // keep dirs at top
			}
			if a.Size != b.Size {
				return a.Size > b.Size
			}
			return strings.ToLower(a.Name) < strings.ToLower(b.Name)
		default: // SortDirsFirst
			if a.IsDir != b.IsDir {
				return a.IsDir
			}
			return strings.ToLower(a.Name) < strings.ToLower(b.Name)
		}
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
