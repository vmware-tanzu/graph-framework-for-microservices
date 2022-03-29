package common

type GraphDBStats struct {
	NodeReadCnt    uint32
	NodeAddCnt     uint32
	NodeDelCnt     uint32
	NodeUpdateCnt  uint32
	LinkReadCnt    uint32
	LinkAddCnt     uint32
	LinkDelCnt     uint32
	LinkUpdateCnt  uint32
	RLinkAddCnt    uint32
	RLinkDelCnt    uint32
	RLinkUpdateCnt uint32
	DbRead         uint32
	DbWrite        uint32
}
