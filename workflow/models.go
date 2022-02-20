package workflow

type TemplateType int
type ArtifactType int

const (
	CONTAINER_SET_TEMPLATE TemplateType = iota
	CONTAINER_TEMPLATE
	DAG_TEMPLATE
	DATA_TEMPLATE
	HTTP_TEMPLATE
	PLUGIN_TEMPLATE
	SCRIPT_TEMPLATE
)

const (
	DEFAULT_ARTIFACT ArtifactType = iota
	GIT_ARTIFACT
	GCS_ARTIFACT
	HDFS_ARTIFACT
	HTTP_ARTICACT
	OSS_ARTIFACT
	RAW_ARTIFACT
	S3_ARTIFACT
)

// TemplateFile represents a single argo workflow template file
type TemplateFile struct {
	Name               string
	Description        string
	Kind               string
	Version            string
	EntrypointTemplate string
	Templates          []*Template
	FilePath           string
	LastUpdatedAt      string
}

// Template represents a single workflow template in the file -> https://argoproj.github.io/argo-workflows/fields/#template
type Template struct {
	Name              string
	Description       string
	Type              TemplateType
	HasInitContainer  bool
	ContainerImageTag string // If Type == CONTAINER_TEMPLATE
	Inputs            *InputOutput
	Outputs           *InputOutput
	Tasks             []*Task // If Type == DAG_TEMPLATE
}

// Task represents a DAGTask in the WorkflowTemplate -> https://argoproj.github.io/argo-workflows/fields/#dagtask
type Task struct {
	Name         string
	Description  string
	Template     string
	Dependencies []string
}

// InputOutput are a common structure to represent inputs and outputs of a template
// inputs: https://argoproj.github.io/argo-workflows/fields/#inputs
// outputs: https://argoproj.github.io/argo-workflows/fields/#outputs
type InputOutput struct {
	Parameters []*Parameter
	Artifacts  []*Artifact
}

// Parameter represents an input or output parameter in the template -> https://argoproj.github.io/argo-workflows/fields/#parameter
type Parameter struct {
	Name        string
	Description string
	Required    bool
}

// Artifact represents an input or output artifact in the template -> https://argoproj.github.io/argo-workflows/fields/#artifact
type Artifact struct {
	Name        string
	Description string
	Required    bool
	Type        ArtifactType
}
