package analyzer

import (
	"os"
	"path/filepath"
	"strings"
)

var targetExtensions = map[string]struct{}{
	".yml":        {},
	".properties": {},
}

func CheckAllJavaFiles(path string) (string, string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return "", ""
	}

	foundAppJava := ""
	foundConfig := ""

	for _, entry := range entries {
		name := entry.Name()
		full := filepath.Join(path, name)

		if entry.IsDir() {
			a, b := CheckAllJavaFiles(full)
			if a != "" {
				foundAppJava = a
			}
			if b != "" {
				foundConfig = b
			}
			continue
		}

		if strings.HasSuffix(name, "Application.java") {
			// colors.Spring.Printf("[ OK ] %s found\n", entry.Name())
			foundAppJava = full
		}

		if _, ok := targetExtensions[filepath.Ext(name)]; ok {
			// colors.Spring.Printf("[ OK ] %s found\n", entry.Name())
			foundConfig = full
		}
	}

	return foundAppJava, foundConfig
}
