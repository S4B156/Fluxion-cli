package analyzer

import "project/models"

func GroupProjects(candidates []models.ProjectCandidate) map[string][]models.ProjectCandidate {
    groups := make(map[string][]models.ProjectCandidate)

    for _, p := range candidates {
        groups[p.ParentFolder] = append(groups[p.ParentFolder], p)
    }

    return groups
}