package workflow

import (
	"errors"
	"gopkg.in/yaml.v3"
)

func parseParameters(node *yaml.Node) ([]*Parameter, error) {
	var result []*Parameter

	for _, parameterNode := range node.Content {
		parameter, err := parseParameter(parameterNode)
		if err != nil {
			return nil, err
		}

		result = append(result, parameter)
	}

	return result, nil
}

func parseParameter(node *yaml.Node) (*Parameter, error) {
	var result Parameter

	if node.Kind == yaml.MappingNode {
		result.Description += cleanupComment(node.HeadComment)
		for index := 0; index < len(node.Content); index += 2 {
			childNode := node.Content[index]
			switch childNode.Value {
			case "name":
				result.Name = node.Content[index+1].Value
				result.Description += "\n"
				result.Description += cleanupComment(node.Content[index+1].LineComment)
			case "value":
				result.Required = false
			case "valueFrom":
				result.Required = false
			}
		}

		return &result, nil
	}

	return nil, errors.New("YAML node is not a mapping node")
}
