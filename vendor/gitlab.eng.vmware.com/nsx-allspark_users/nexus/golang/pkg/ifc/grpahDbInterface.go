package ifc

import (
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/common"
)

type GraphDBInterface interface {
	// un-anchored nodes will need the nodeName to be globally Unique.
	UpsertNode(requestorId, nodeType, nodeName string, nodeProp PropertyType) *GNode
	// create a node that is linked to a parent node,
	// under the parent the "linkKey" field is used for uniquely identify the node.
	UpsertChildNode(requestorId, parentNodeId, linkType, linkKey string,
		linkProp PropertyType,
		nodeType string,
		nodeProp PropertyType) (*GNode, *GLink)

	DescribeNode(nodeId, nodeType string) *GNode
	DeleteNode(nodeId string)
	UpdateNodeAddProperties(requestorId, nodeId string, properties PropertyType) int64
	UpdateNodeRemoveProperties(requestorId, nodeId string, properties []string) int64

	UpsertLink(requestorId string, lnkObj common.UpsertLinkOpObj) (*GLink, *GLink)
	/*	UpsertLink(
		requestorId, linkType, sourceNodeId, destinationNodeId, sourceNodeType, destinationNodeType string,
		sourceNodeKeyValue, destinationNodeKeyValue string,
		sourceNodeProperties, destinationNodeProperties, linkProperties PropertyType,
		isSingletonLink bool) (*GLink, *GLink)
	*/
	DescribeLinkById(linkId string) (*GLink, bool)
	DescribeLink(linkType, sourceId, destinaitonId string) (*GLink, bool)
	DeleteLink(linkId string)
	DeleteSoftLinkWithRSoftLink(linkId string, rLinkId string)
	DeleteLinks(nodeId string) []*GLink
	UpdateLinkAddProperties(
		requestorId, linkId string, properties PropertyType) int64
	UpdateLinkRemoveProperties(requestorId, linkId string, properties []string) int64
	// iterator
	GetNextLinkKey(sourceId, linkType, currentKey string, cnt int) []string
	// rlink iterator
	GetNextRLinkKey(sourceId, linkType, currentKey string, cnt int) []string
	// just upate node property
	GetNodeProperty(nodeId string) PropertyType
	GetStats() common.GraphDBStats
	Shutdown()
	GetLatestRevision() int64
	CompactRevision(int64, int64, bool) (*EtcdCompactResponse, error)
}
