package workflow

import (
	"errors"

	"gopkg.in/yaml.v3"
)

func parseArtifacts(node *yaml.Node) ([]*Artifact, error) {
	var result []*Artifact

	for _, parameterNode := range node.Content {
		parameter, err := parseArtifact(parameterNode)
		if err != nil {
			return nil, err
		}

		result = append(result, parameter)
	}

	return result, nil
}

func parseArtifact(node *yaml.Node) (*Artifact, error) {
	var result Artifact

	if node.Kind == yaml.MappingNode {
		result.Description += cleanupComment(node.HeadComment)
		for index := 0; index < len(node.Content); index += 2 {
			childNode := node.Content[index]
			switch childNode.Value {
			case "name":
				result.Name = node.Content[index+1].Value
				result.Description += "\n"
				result.Description += cleanupComment(node.Content[index+1].LineComment)
			case "gcs":
				result.Type = GCS_ARTIFACT
				result.Required = false
			case "git":
				result.Type = GIT_ARTIFACT
				result.Required = false
			case "hdfs":
				result.Type = HDFS_ARTIFACT
				result.Required = false
			case "http":
				result.Type = HTTP_ARTICACT
				result.Required = false
			case "oss":
				result.Type = OSS_ARTIFACT
				result.Required = false
			case "raw":
				result.Type = RAW_ARTIFACT
				result.Required = false
			case "s3":
				result.Type = S3_ARTIFACT
				result.Required = false
			}
		}

		return &result, nil
	}

	return nil, errors.New("YAML node is not a mapping node")
}
