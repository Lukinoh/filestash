package utils

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

var MOCK_CURRENT_DIR string

func GetCurrentDir() string {
	if MOCK_CURRENT_DIR != "" {
		return MOCK_CURRENT_DIR
	}
	ex, _ := os.Executable()
	return filepath.Dir(ex)
}

func GetAbsolutePath(base string, opts ...string) string {
	fullPath := base
	if strings.HasPrefix(base, "/") == false { // relative filepath are relative to the binary
		fullPath = filepath.Join(GetCurrentDir(), base)
	}
	if len(opts) == 0 {
		return fullPath
	}
	return filepath.Join(append([]string{fullPath}, opts...)...)
}

func IsDirectory(path string) bool {
	if path == "" {
		return false
	}
	if path[len(path)-1:] != "/" {
		return false
	}
	return true
}

/*
 * Join 2 path together, result has a file
 */
func JoinPath(base, file string) string {
	filePath := path.Join(base, file)
	if strings.HasPrefix(filePath, base) == false {
		return base
	}
	return filePath
}

func EnforceDirectory(path string) string {
	if path == "" {
		return "/"
	} else if path[len(path)-1:] == "/" {
		return path
	}
	return path + "/"
}

func SplitPath(path string) (root string, filename string) {
	if path == "" {
		path = "/"
	}
	if IsDirectory(path) == false {
		filename = filepath.Base(path)
	}
	if root = strings.TrimSuffix(path, filename); root == "" {
		root = "/"
	}
	return root, filename
}

func GlobMatch(pattern, name string) bool {
	m, _ := doublestar.Match(pattern, name)
	return m
}

/*
 * Normalise the Path in the form of "/xxx/xxx" or simply "/"
 */
func NormalisePath(path string) string {
	path = strings.TrimSpace(path)
	path = strings.Trim(path, "/")
	return "/" + path
}

/*
 * Normalised a comma separated list of path into an array:
 * ParsePathList("docs")            → ["/docs"]
 * ParsePathList("/docs/")          → ["/docs"]
 * ParsePathList("docs, other")     → ["/docs", "/other"]
 * ParsePathList("docs,,/secret/")  → ["/docs", "/secret"]
 */
func ParsePathList(raw string) []string {
	list := make([]string, 0)
	for _, path := range strings.Split(raw, ",") {
		path = NormalisePath(path)

		if path != "/" {
			list = append(list, path)
		}
	}
	return list
}

/*
 * Express "fullpath" relative to "chroot", always starting with "/"
 * RelativePath("/data", "/data/docs/secret")  → "/docs/secret"
 * RelativePath("/data", "/data")              → "/"
 * RelativePath("/data", "/data/secret")       → "/secret"
 */
func RelativePath(chroot, fullpath string) string {
	return NormalisePath(strings.TrimPrefix(fullpath, strings.TrimSuffix(EnforceDirectory(chroot), "/")))
}

/*
 * Check if the relPath is part of the denyList.
 * Either by exact matching or if it is nested
 */
func IsPathDenied(denyList []string, relPath string) bool {
	relPath = NormalisePath(relPath)
	for _, entry := range denyList {
		if relPath == entry || strings.HasPrefix(relPath, entry+"/") {
			return true
		}
	}
	return false
}

/*
 * Remove the directory entries that fall under a deny-listed path
 */
func FilterPathDenyList(denyList []string, dirRelPath string, entries []os.FileInfo) []os.FileInfo {
	if len(denyList) == 0 {
		return entries
	}

	dirRelPath = strings.TrimSuffix(dirRelPath, "/")
	filtered := make([]os.FileInfo, 0, len(entries))
	for _, entry := range entries {
		if !IsPathDenied(denyList, dirRelPath+"/"+entry.Name()) {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}
