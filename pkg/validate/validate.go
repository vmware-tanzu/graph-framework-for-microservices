package validate

var CrdParentsMap = map[string][]string{}

type NexusAnnotation struct {
	Hierarchy []string `json:"hierarchy"`
}
