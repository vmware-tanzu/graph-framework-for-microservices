package nexus

const (
	tagName              = "nexus.vmware"
	BaseGroupName string = "tanzu.tsm.vmware.com"
)

type ID struct {
	Id string `nexus.vmware:"id"`
}

type Node struct {
	ID
}

type RestUris struct {
	Uri     string   `json:"uri"`
	Methods []string `json:"methods"`
}

type RestAPISpec struct {
	Uris []RestUris `json:"uris"`
}
