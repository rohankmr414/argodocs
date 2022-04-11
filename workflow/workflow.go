package workflow

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
	"time"
)

// ParseFiles accepts a glob pattern for argo workflow template files, parses them and returns them as TemplateFile structs
func ParseFiles(pattern string) ([]*TemplateFile, error) {
	var result []*TemplateFile

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	for _, path := range matches {
		var templateFile *TemplateFile
		templateFile, err = parseFile(path)
		if err != nil {
			continue
		}
		result = append(result, templateFile)
	}

	return result, nil
}

// parseFile parses a single argo workflow template file and returns that as a TemplateFile object
func parseFile(path string) (*TemplateFile, error) {
	yamlNode := yaml.Node{}

	yamlFileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFileContent, &yamlNode)
	if err != nil {
		return nil, err
	}

	var templateFile *TemplateFile
	templateFile, err = parseTemplateFile(&yamlNode)
	if err != nil {
		return nil, err
	}

	templateFile.FilePath = path
	templateFile.LastUpdatedAt = time.Now().Format(time.RFC850)

	if templateFile.EntrypointTemplate == "" {
		templateFile.EntrypointTemplate = "nil"
	}
	
	return templateFile, nil
}
