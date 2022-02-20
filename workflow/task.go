package workflow

import (
	"errors"
	"gopkg.in/yaml.v3"
	"strings"
)

func parseTasks(node *yaml.Node) ([]*Task, error) {
	var result []*Task
	if node.Kind == yaml.MappingNode {
		for index := 0; index < len(node.Content); index += 2 {
			childNode := node.Content[index]
			switch childNode.Value {
			case "tasks":
				taskNodes := node.Content[index+1].Content
				for _, taskNode := range taskNodes {
					task, err := parseTask(taskNode)
					if err != nil {
						return nil, err
					}

					result = append(result, task)
				}
			}
		}

		return result, nil
	}

	return nil, errors.New("YAML node is not a mapping node")
}

func parseTask(node *yaml.Node) (*Task, error) {
	var result Task

	if node.Kind == yaml.MappingNode {
		result.Description += cleanupComment(node.HeadComment)
		for index := 0; index < len(node.Content); index += 2 {
			childNode := node.Content[index]
			switch childNode.Value {
			case "name":
				result.Name = node.Content[index+1].Value
				result.Description += "\n"
				result.Description += cleanupComment(node.Content[index+1].LineComment)
			case "template":
				result.Template = node.Content[index+1].Value
			case "templateRef":
				name, templateName, err := parseTemplateRef(node.Content[index+1])
				if err != nil {
					return nil, err
				}
				result.Template = strings.Join([]string{name, templateName}, "/")
			case "dependencies":
				var dependencies []string
				for _, dependencyNode := range node.Content[index+1].Content {
					dependencies = append(dependencies, dependencyNode.Value)
				}
				result.Dependencies = dependencies
			}
		}

		return &result, nil
	}

	return nil, errors.New("YAML node is not a mapping node")
}

func parseTemplateRef(node *yaml.Node) (string, string, error) {
	var name string
	var template string

	if node.Kind == yaml.MappingNode {
		for index := 0; index < len(node.Content); index += 2 {
			childNode := node.Content[index]
			switch childNode.Value {
			case "name":
				name = node.Content[index+1].Value
			case "template":
				template = node.Content[index+1].Value
			}
		}

		return name, template, nil
	}

	return name, template, errors.New("YAML node is not a mapping node")
}
