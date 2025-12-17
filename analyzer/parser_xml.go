package analyzer

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"project/models"
	"strings"
)

func AnalyzeProject(pathToPom string) (*models.ProjectCandidate, error) {
	contentBytes, err := os.ReadFile(pathToPom)
	if err != nil {
		return nil, err
	}

	rootPath := filepath.Dir(pathToPom)
	isSpring := strings.Contains(string(contentBytes), "spring-boot")

	if isSpring {
		candidate := new(models.ProjectCandidate)
		dir := filepath.Join(rootPath, "src")
		foundAppJava, foundConfig := CheckAllJavaFiles(dir)

		if foundAppJava == "" || foundConfig == "" {
			return nil, nil
		}

		// --- Initializing the config file ---
		candidate.ApplicationFilePath = foundAppJava
		
		ext := filepath.Ext(foundConfig)
		if ext == ".yml" || ext == ".yaml" {
			candidate.Config = ParsingYaml(foundConfig)
		} else {
			candidate.Config, err = ParsingProperties(foundAppJava)
			if err != nil {
				return nil, err
			}
		}
		// --- end ---

		candidate, err = parsingXmlDependencies(pathToPom, candidate)
		if err != nil {
			return nil, err
		}

		// --- Initializing the paths ---
		parent := filepath.Dir(rootPath)
		candidate.Path = rootPath
		candidate.ParentFolder = parent
		// --- end ---
		
		return candidate, nil
	}
	return nil, nil
}

func parsingXmlDependencies(path string, candidate *models.ProjectCandidate) (*models.ProjectCandidate, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var project models.Project
    err = xml.NewDecoder(file).Decode(&project)
    if err != nil {
        return nil, err
    }

	// --- Initializing Spring Features ---
	features := models.ServiceFeatures{}

	for _, dep := range project.Dependencies.Dependency {
		art := dep.ArtifactID
		if strings.Contains(art, "postgresql") {
			features.HasPostgres = true
		}
		if strings.Contains(art, "mysql") {
			features.HasMySQL = true
		}
		if strings.Contains(art, "redis") || strings.Contains(art, "jedis") || strings.Contains(art, "lettuce") {
			features.HasRedis = true
		}
		if strings.Contains(art, "consul") {
			features.HasConsul = true
		}
		if strings.Contains(art, "kafka") {
			features.HasKafka = true
		}
		if strings.Contains(art, "rabbit") {
			features.HasKafka = true
		}
		if strings.Contains(art, "eureka") {
			features.HasEureka = true
		}
	}

	candidate.ServiceFeatures = features
	// --- end ---

    candidate.Name = project.Name
	candidate.MetaData = project
	if project.Dependencies.Dependency != nil {
        for _, dep := range project.Dependencies.Dependency {
            candidate.Dependencies = append(candidate.Dependencies, dep.ArtifactID)
        }
    }
    return candidate, nil
}