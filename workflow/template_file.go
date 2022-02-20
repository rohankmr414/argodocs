package workflow

import "gopkg.in/yaml.v3"

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
