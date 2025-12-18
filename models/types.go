package models

import (
	"encoding/xml"
	"strconv"
)

type Project struct {
	Parent struct {
		SpringVersion string `xml:"version"`
	}
	XMLName     xml.Name `xml:"project"`
	GroupID     string   `xml:"groupId"`
	ArtifactID  string   `xml:"artifactId"`
	Name        string   `xml:"name"`
	Description string   `xml:"description"`
	Properties  struct {
		JavaVersion        string `xml:"java.version"`
		SpringCloudVersion string `xml:"spring-cloud.version"`
	} `xml:"properties"`
	Dependencies struct {
		Dependency []struct {
			GroupID    string `xml:"groupId"`
			ArtifactID string `xml:"artifactId"`
		} `xml:"dependency"`
	} `xml:"dependencies"`
}

type AppService struct {
	Name     string
	Path     string
	Port     int
	Features ServiceFeatures
}

type GlobalConfig struct {
	Version  string
	Services []AppService
	Features ServiceFeatures
	DbName   string
	DbUser   string
	DbPass   string
}

type ProjectCandidate struct {
	Path                string
	Config              any
	ApplicationFilePath string
	MetaData            Project
	ServiceFeatures     ServiceFeatures
	Name                string
	Dependencies        []string
	ParentFolder        string
}

type ServiceFeatures struct {
	// Databases
	HasPostgres  bool
	HasMySQL     bool
	HasMariaDB   bool
	HasMongo     bool
	HasCassandra bool
	HasRedis     bool
	HasElastic   bool

	// Messaging
	HasKafka    bool
	HasRabbit   bool
	HasActiveMQ bool

	// Spring Cloud / Infrastructure
	HasEureka       bool
	HasConsul       bool
	HasConfigClient bool
	HasConfigServer bool
	HasGateway      bool
	HasFeign        bool
	
	// Observability & Security
	HasZipkin     bool
	HasPrometheus bool
	HasVault      bool
}

type ComposeData struct {
	ServiceName string
	AppPort     int
	Features    ServiceFeatures

	DbName string
	DbUser string
	DbPass string
}

// func (p ProjectCandidate) NeedCompose() bool {
// 	features := p.ServiceFeatures
// 	return features.HasPostgres || features.HasMySQL || features.HasRedis || features.HasMongo ||
// 		features.HasKafka || features.HasRabbit || features.HasEureka
// }

func (p ProjectCandidate) GetPort() int {
	defaultPort := 8080

	if p.Config == nil {
		return defaultPort
	}

	switch v := p.Config.(type) {
	case Config:
		if v.Server.Port != 0 {
			return v.Server.Port
		}
	case map[string]string:
		val, ok := v["server.port"]
		if ok {
			port, err := strconv.Atoi(val)
			if err == nil {
				return port
			}
		}
	}
	return defaultPort
}

type Config struct {
	Spring struct {
		Application struct {
			Name string `yaml:"name"`
		} `yaml:"application"`
		Datasource struct {
			URL             string `yaml:"url"`
			Username        string `yaml:"username"`
			Password        string `yaml:"password"`
			DriverClassName string `yaml:"driver-class-name"`
		} `yaml:"datasource"`
		Jpa struct {
			Database string `yaml:"database"`
		} `yaml:"jpa"`
	} `yaml:"spring"`
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
}
