package analyzer

import (
	"project/models"
	"strings"
)

func IsMicroserviceSystem(projects []models.ProjectCandidate) bool {
	for _, p := range projects {
		for _, dep := range p.Dependencies {
			if strings.Contains(dep, "spring-cloud-starter-gateway") ||
			   strings.Contains(dep, "spring-cloud-starter-netflix-eureka-server") ||
			   strings.Contains(dep, "spring-cloud-config-server") ||
			   strings.Contains(dep, "spring-cloud-starter-openfeign") ||
			   strings.Contains(dep, "spring-cloud-starter-circuitbreaker") ||
			   strings.Contains(dep, "spring-cloud-starter-sleuth") ||
			   strings.Contains(dep, "spring-cloud-starter-consul") ||
			   strings.Contains(dep, "spring-cloud-stream") {
				return true
			}
		}
		// Additional check: if there are many projects and they have a common parent pom (not implemented yet),
		// but the Gateway check is the most reliable.
	}
	return false
}
