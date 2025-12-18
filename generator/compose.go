package generator

import (
	_ "embed"
	"os"
	"path/filepath"
	"project/models"
	"project/pkg/colors"
	"text/template"
)

//go:embed templates/docker-compose.tpl
var dockerComposeTpl string

func GenerateGlobalCompose(rootPath string, candidates []models.ProjectCandidate) {
	var services []models.AppService
	var globalFeatures models.ServiceFeatures

	for _, c := range candidates {
		relativePath, _ := filepath.Rel(rootPath, c.Path)
		
		svc := models.AppService{
			Name:     c.Name,
			Path:     "./" + filepath.ToSlash(relativePath),
			Port:     c.GetPort(),
			Features: c.ServiceFeatures,
		}
		services = append(services, svc)

		// A bit ugly, but robust aggregation of infrastructure needs
		if c.ServiceFeatures.HasPostgres { globalFeatures.HasPostgres = true }
		if c.ServiceFeatures.HasMySQL { globalFeatures.HasMySQL = true }
		if c.ServiceFeatures.HasMariaDB { globalFeatures.HasMariaDB = true }
		if c.ServiceFeatures.HasRedis { globalFeatures.HasRedis = true }
		if c.ServiceFeatures.HasMongo { globalFeatures.HasMongo = true }
		if c.ServiceFeatures.HasCassandra { globalFeatures.HasCassandra = true }
		if c.ServiceFeatures.HasElastic { globalFeatures.HasElastic = true }
		
		if c.ServiceFeatures.HasKafka { globalFeatures.HasKafka = true }
		if c.ServiceFeatures.HasRabbit { globalFeatures.HasRabbit = true }
		if c.ServiceFeatures.HasActiveMQ { globalFeatures.HasActiveMQ = true }

		if c.ServiceFeatures.HasEureka { globalFeatures.HasEureka = true }
		if c.ServiceFeatures.HasConsul { globalFeatures.HasConsul = true }
		
		// Infrastructure for Observability / Security
		if c.ServiceFeatures.HasZipkin { globalFeatures.HasZipkin = true }
		if c.ServiceFeatures.HasVault { globalFeatures.HasVault = true }
		
		// Note: Gateway and ConfigServer are usually apps within candidates, 
		// so we don't necessarily spin up "external" containers for them here,
		// unless we want to use standard images instead of custom code.
	}

	config := models.GlobalConfig{
		Version:  "3.9",
		Services: services,
		Features: globalFeatures,
		DbName:   "ghost_main_db", 
		DbUser:   "admin",
		DbPass:   "secret",
	}

	writeComposeFile(rootPath, config)
}

func GenerateSingleCompose(candidate models.ProjectCandidate) {
	service := models.AppService{
		Name:     candidate.Name,
		Path:     ".",
		Port:     candidate.GetPort(),
		Features: candidate.ServiceFeatures,
	}

	config := models.GlobalConfig{
		Version:  "3.9",
		Services: []models.AppService{service},
		Features: candidate.ServiceFeatures,
		DbName:   "ghost_single_db",
		DbUser:   "user",
		DbPass:   "password",
	}

	writeComposeFile(candidate.Path, config)
}

func writeComposeFile(path string, data models.GlobalConfig) {
	tmpl, err := template.New("docker-compose").Parse(dockerComposeTpl)
	if err != nil {
		colors.Error.Printf("Error parsing template: %v\n", err)
		return
	}

	filePath := filepath.Join(path, "docker-compose.yml")
	file, err := os.Create(filePath)
	if err != nil {
		colors.Error.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	err = tmpl.Execute(file, data)
	if err != nil {
		colors.Error.Printf("Error executing template: %v\n", err)
		return
	}
	colors.Docker.Printf("[ DONE ] Docker Compose created at %s\n", filePath)
}
