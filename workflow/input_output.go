package workflow

import (
	"errors"
	"gopkg.in/yaml.v3"
)

func parseInputsOutputs(node *yaml.Node) (*InputOutput, error) {
	var result InputOutput

	if node.Kind == yaml.MappingNode {
		for index := 0; index < len(node.Content); index += 2 {
			childNode := node.Content[index]
			if childNode.Kind == yaml.ScalarNode {
				switch childNode.Value {
				case "parameters":
					parameters, err := parseParameters(node.Content[index+1])
					if err != nil {
						return nil, err
					}
					result.Parameters = parameters
				case "artifacts":
					artifacts, err := parseArtifacts(node.Content[index+1])
					if err != nil {
						return nil, err
					}
					result.Artifacts = artifacts
				}
			}
		}

		return &result, nil
	}

	return nil, errors.New("YAML node is not a mapping node")
}
