package ifc

type GNode struct {
	Id         string
	Type       string
	Properties PropertyType
	Links      []*GLink
	RLinks     []*GLink
}

type EtcdCompactResponse struct {
	ClusterId uint64 `json:"cluster-id"`
	MemberId  uint64 `json:"member-id"`
	RaftTerm  uint64 `json:"raft-term"`
	Revision  int64  `json:"revision"`
}
