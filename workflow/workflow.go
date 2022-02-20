package workflow

import (
	"errors"
	"github.com/junaidrahim/argodocs/logger"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
	"strings"
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

func parseTemplateFile(node *yaml.Node) (*TemplateFile, error) {
	var result TemplateFile
	node = node.Content[0] // root map node

	if node.Kind == yaml.MappingNode {
		for index := 0; index < len(node.Content); index += 2 {
			child := node.Content[index]
			if child.Kind == yaml.ScalarNode {
				switch child.Value {
				case "apiVersion":
					result.Version = node.Content[index+1].Value
					result.Description += cleanupComment(node.Content[index].HeadComment)
				case "kind":
					result.Kind = node.Content[index+1].Value
					result.Description += cleanupComment(node.Content[index].HeadComment)
				case "metadata":
					metadataMapNode := node.Content[index+1]
					if metadataMapNode.Kind == yaml.MappingNode {
						for metadataChildIndex := 0; metadataChildIndex < len(metadataMapNode.Content); metadataChildIndex += 2 {
							metadataChildNode := metadataMapNode.Content[metadataChildIndex]
							if metadataChildNode.Value == "name" || metadataChildNode.Value == "generateName" {
								result.Name = metadataMapNode.Content[metadataChildIndex+1].Value
							}
						}
					}
				case "spec":
					specMapNode := node.Content[index+1]
					if specMapNode.Kind == yaml.MappingNode {
						for specChildIndex := 0; specChildIndex < len(specMapNode.Content); specChildIndex += 2 {
							specChildNode := specMapNode.Content[specChildIndex]
							if specChildNode.Kind == yaml.ScalarNode {
								if specChildNode.Value == "entrypoint" {
									result.EntrypointTemplate = specMapNode.Content[specChildIndex+1].Value
								} else if specChildNode.Value == "templates" {
									// this is the template mapping node, need to pass this off to another func
									templates, _ := parseTemplates(specMapNode.Content[specChildIndex+1])
									result.Templates = templates
								}
							}
						}
					}
				}
			}
		}
	}

	return &result, nil
}

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
					// add container image tag
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
