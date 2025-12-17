package analyzer

import (
	"bufio"
	"os"
	"strings"
	"project/models"
	"gopkg.in/yaml.v3"
)

func ParsingYaml(path string) models.Config {
	file, err := os.Open(path)
	if err != nil {
		return models.Config{}
	}
	defer file.Close()

	var config models.Config
	err = yaml.NewDecoder(file).Decode(&config)
	if err != nil {
		return models.Config{}
	}
	return config
}

func ParsingProperties(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "!") {
			continue
		}

		equal := strings.Index(line, "=")
		if equal < 0 {
			continue
		}

		key := strings.TrimSpace(line[:equal])
		value := ""
		if len(line) > equal+1 {
			value = strings.TrimSpace(line[equal+1:])
		}

		if len(key) > 0 {
			config[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return config, nil
}
