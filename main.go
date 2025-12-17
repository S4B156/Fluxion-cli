package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"project/analyzer"
	"project/generator"
	"project/models"
	"project/pkg/colors"
)

func main() {
	// --- Phase 0: Getting the path ---
	pathPtr := flag.String("path", "", "path to the project for analysis")
	flag.Parse()

	path := *pathPtr
	if path == "" {
		log.Fatal("Error: project path not provided (use -path flag)")
	}
	path = filepath.Clean(path)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("Error: directory does not exist: %s", path)
	}
	// --- end ---

	// --- Phase 1: Data Collection (Scanner) ---
	projects := []models.ProjectCandidate{}
	projects = analyzer.ScanAllProjects(path)
	colors.Info.Println(len(projects), "projects found")
	// --- end ---

	// --- Phase 2: Grouping and Analysis (The Brain) ---
	groups := analyzer.GroupProjects(projects)

	colors.Info.Printf("Found %d project groups\n", len(groups))

	for parentPath, groupMembers := range groups {
		isMicroservices := analyzer.IsMicroserviceSystem(groupMembers)
		// --- Phase 3: Generation ---
		if isMicroservices {
			colors.Detected.Printf("Detected Microservice System at: %s\n", parentPath)
			generator.GenerateGlobalCompose(parentPath, groupMembers)

			for _, proj := range groupMembers {
				generator.GenerateDockerfile(proj.Path, proj.MetaData, proj.GetPort())
			}

		} else {
			colors.Info.Printf("Detected Standalone Projects at: %s\n", parentPath)

			for _, proj := range groupMembers {
				port := proj.GetPort()

				generator.GenerateDockerfile(proj.Path, proj.MetaData, port)

				generator.GenerateSingleCompose(proj)
			}
		}
		// --- end ---
	}
	// --- end ---
}
