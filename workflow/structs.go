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

type Task struct {
	Name         string
	Description  string
	Template     string
	Dependencies []string
}

type InputOutput struct {
	Parameters []*Parameter
	Artifacts  []*Artifact
}

type Parameter struct {
	Name        string
	Description string
	Required    bool
}

type Artifact struct {
	Name        string
	Description string
	Required    bool
	Type        ArtifactType
}
