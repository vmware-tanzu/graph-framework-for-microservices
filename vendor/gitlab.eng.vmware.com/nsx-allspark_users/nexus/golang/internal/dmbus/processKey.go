package dmbus

import (
	"strings"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/internal/graphdb"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/ifc"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/logging"
)

type ProcessKey struct {
	KeyList                  []string
	Created                  bool
	RLink                    bool
	Link                     bool
	Lock                     bool
	UnknownNodeType          bool
	Path                     ifc.NodePathList
	NodeType                 string
	NodeKey                  string
	ParentId                 string
	NodeId                   string
	LinkInfoDestNodeType     string
	LinkInfoDestNodeKeyValue string
}

func NewProcessKey() *ProcessKey {
	e := &ProcessKey{
		Link:                     false,
		RLink:                    false,
		Lock:                     false,
		Created:                  false,
		UnknownNodeType:          false,
		KeyList:                  []string{},
		Path:                     ifc.NodePathList{},
		LinkInfoDestNodeType:     "",
		LinkInfoDestNodeKeyValue: "",
	}
	return e
}
func (p *ProcessKey) process(k string) *ProcessKey {
	p.KeyList = strings.Split(k, "/")
	p.Path = ifc.NodePathList{}
	p.Lock = false
	p.Created = false
	p.Link = false
	p.RLink = false
	p.UnknownNodeType = false
	p.NodeType = ""
	p.NodeKey = ""
	validObjectTypes := map[string]struct{}{graphdb.FixedPath_Links: struct{}{},
		graphdb.FixedPath_RLinks: struct{}{}, graphdb.FixedPath_Created: struct{}{},
		graphdb.FixedPath_Lock: struct{}{}}

	for i := 1; i < len(p.KeyList); i = i + 2 {
		ntype := p.KeyList[i]
		if ntype[0] == '_' {
			if _, ok := validObjectTypes[ntype]; !ok {
				p.UnknownNodeType = true
				logging.Debugf("Encountered unknown node object type"+
					" : %s", ntype)
				break
			}
		}
		if ntype == graphdb.FixedPath_Created {
			p.Created = true
			break
		} else if ntype == graphdb.FixedPath_Links || ntype == graphdb.FixedPath_RLinks {
			if ntype == graphdb.FixedPath_Links {
				p.Link = true
			} else if ntype == graphdb.FixedPath_RLinks {
				p.RLink = true
			}
			p.LinkInfoDestNodeType = p.KeyList[i+1]
			p.LinkInfoDestNodeKeyValue = p.KeyList[i+2]
			if len(p.KeyList) > (i + 3) {
				p.Created = p.KeyList[i+3] == graphdb.FixedPath_Created
				p.Lock = p.KeyList[i+3] == graphdb.FixedPath_Lock
			}
			break
		} else if ntype == graphdb.FixedPath_Lock {
			p.Lock = true
			break
		}
		nodeVal := p.KeyList[i+1]
		p.NodeType = ntype
		p.NodeKey = nodeVal

		p.Path = append(p.Path, ifc.NodePath{ntype, nodeVal})
	}
	p.ParentId = ""
	p.NodeId = ""
	prevPath := ""
	if p.UnknownNodeType == false {
		for _, np := range p.Path {
			tmpVar := "/" + np[0] + "/" + np[1]
			p.NodeId += tmpVar
			p.ParentId += prevPath
			prevPath = tmpVar
		}
	}
	return p
}
