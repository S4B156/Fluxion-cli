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
			candidate.Config, err = ParsingProperties(foundConfig)
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
		art := strings.ToLower(dep.ArtifactID)

		// Databases
		if strings.Contains(art, "postgresql") {
			features.HasPostgres = true
		}
		if strings.Contains(art, "mysql") {
			features.HasMySQL = true
		}
		if strings.Contains(art, "mariadb") {
			features.HasMariaDB = true
		}
		if strings.Contains(art, "redis") || strings.Contains(art, "jedis") || strings.Contains(art, "lettuce") {
			features.HasRedis = true
		}
		if strings.Contains(art, "mongodb") {
			features.HasMongo = true
		}
		if strings.Contains(art, "cassandra") {
			features.HasCassandra = true
		}
		if strings.Contains(art, "elasticsearch") || strings.Contains(art, "opensearch") {
			features.HasElastic = true
		}

		// Messaging
		if strings.Contains(art, "kafka") {
			features.HasKafka = true
		}
		if strings.Contains(art, "rabbit") || strings.Contains(art, "amqp") {
			features.HasRabbit = true
		}
		if strings.Contains(art, "activemq") || strings.Contains(art, "artemis") {
			features.HasActiveMQ = true
		}

		// Spring Cloud & Discovery
		if strings.Contains(art, "eureka") {
			features.HasEureka = true
		}
		if strings.Contains(art, "consul") {
			features.HasConsul = true
		}
		if strings.Contains(art, "spring-cloud-starter-config") {
			features.HasConfigClient = true
		}
		if strings.Contains(art, "spring-cloud-config-server") {
			features.HasConfigServer = true
		}
		if strings.Contains(art, "spring-cloud-starter-gateway") {
			features.HasGateway = true
		}
		if strings.Contains(art, "feign") {
			features.HasFeign = true
		}

		// Observability & Security
		if strings.Contains(art, "zipkin") || strings.Contains(art, "sleuth") || strings.Contains(art, "micrometer-tracing") {
			features.HasZipkin = true
		}
		if strings.Contains(art, "prometheus") || strings.Contains(art, "micrometer") {
			features.HasPrometheus = true
		}
		if strings.Contains(art, "vault") {
			features.HasVault = true
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