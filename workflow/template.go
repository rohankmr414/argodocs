package workflow

import (
	"errors"
	"gopkg.in/yaml.v3"
)

func parseTemplates(node *yaml.Node) ([]*Template, error) {
	var result []*Template
	if node.Kind == yaml.SequenceNode {
		for _, templateNode := range node.Content {
			template, err := parseTemplate(templateNode)
			if err != nil {
				return nil, err
			}

			result = append(result, template)
		}

		return result, nil
	}

	return nil, errors.New("YAML node is not a sequence node")
}

func parseTemplate(node *yaml.Node) (*Template, error) {
	var result Template
	if node.Kind == yaml.MappingNode {
		result.Description += cleanupComment(node.HeadComment)
		for index := 0; index < len(node.Content); index += 2 {
			childNode := node.Content[index]
			if childNode.Kind == yaml.ScalarNode {
				switch childNode.Value {
				case "name":
					result.Name = node.Content[index+1].Value
					result.LineNumber = childNode.Line
				case "inputs":
					inputs, err := parseInputsOutputs(node.Content[index+1])
					if err != nil {
						return nil, err
					}
					result.Inputs = inputs
				case "outputs":
					outputs, err := parseInputsOutputs(node.Content[index+1])
					if err != nil {
						return nil, err
					}
					result.Outputs = outputs
				case "containerSet":
					result.Type = CONTAINER_SET_TEMPLATE
				case "container":
					result.Type = CONTAINER_TEMPLATE
					containerImage, err := parseContainerImage(node.Content[index+1])
					if err != nil {
						return nil, err
					}
					result.ContainerImageTag = containerImage
				case "dag":
					result.Type = DAG_TEMPLATE
					tasks, err := parseTasks(node.Content[index+1])
					if err != nil {
						return nil, err
					}
					result.Tasks = tasks
				case "data":
					result.Type = DATA_TEMPLATE
				case "http":
					result.Type = HTTP_TEMPLATE
				case "plugin":
					result.Type = PLUGIN_TEMPLATE
				case "script":
					result.Type = SCRIPT_TEMPLATE
				case "initContainers":
					result.HasInitContainer = true
				}
			}
		}

		return &result, nil
	}

	return nil, errors.New("YAML node is not a mapping node")
}

func parseContainerImage(node *yaml.Node) (string, error) {
	if node.Kind == yaml.MappingNode {
		for index := 0; index < len(node.Content); index += 2 {
			childNode := node.Content[index]
			switch childNode.Value {
			case "image":
				containerTag := node.Content[index+1].Value
				return containerTag, nil
			}
		}
	}

	return "", errors.New("YAML node is not a mapping node")
}
