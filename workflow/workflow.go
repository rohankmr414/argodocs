package workflow

import (
	"github.com/junaidrahim/argodocs/logger"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
	"time"
)

var l = logger.GetLogger("Workflow: ")

func ParseFiles(pattern string) ([]*TemplateFile, error) {
	l.Printf("hello %v\n", pattern)
	var result []*TemplateFile

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	for _, path := range matches {
		l.Printf(path)
		var templateFile *TemplateFile
		templateFile, err = parseFile(path)
		if err != nil {
			l.Printf("error for parsing file %v", path)
			continue
		}
		result = append(result, templateFile)
	}

	return result, nil
}

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
	return templateFile, nil
}
