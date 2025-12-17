package analyzer

import (
	"os"
	"path/filepath"
	"project/models"
)

func ScanAllProjects(rootPath string) []models.ProjectCandidate {
    var candidates []models.ProjectCandidate
    scanRecursive(rootPath, &candidates)
    return candidates
}

func scanRecursive(path string, candidates *[]models.ProjectCandidate) {
    entries, err := os.ReadDir(path)
    if err != nil {
        return
    }

    for _, entry := range entries {
        fullPath := filepath.Join(path, entry.Name())

        if entry.IsDir() {
            scanRecursive(fullPath, candidates)
        } else if entry.Name() == "pom.xml" {
            candidate, err := AnalyzeProject(fullPath)
            if err == nil && candidate != nil {
                *candidates = append(*candidates, *candidate)
            }
        }
    }
}
