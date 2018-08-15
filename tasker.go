package tasker

// A Tasker performs special operations on ECS Task Definitions.
type Tasker interface {
	UpdateContainerDefinition(in UpdateContainerInput) (string, error)
	UpdateContainerImage(in ImageUpdateInput) error
}

type UpdateContainerInput struct {
	Cluster       string
	Service       string
	Family        string
	ContainerDefs string
}

type ImageUpdateInput struct {
	Cluster string
	Service string
	Family  string
	Image   string
}
