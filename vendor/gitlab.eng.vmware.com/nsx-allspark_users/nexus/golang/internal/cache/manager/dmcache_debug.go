package manager

import (
	"fmt"

	. "gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/internal/cache/subnode"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/base"
)

type SubTreeDebugShow struct {
	*SubNodeType
	CbfnNodeCount int                         `json:"cbfnNodeCount"`
	CbfnLinkCount int                         `json:"cbfnLinkCount"`
	Child         map[string]SubTreeDebugShow `json:"child"`
}

func (c *CacheManagerSubTreeMap) setSubtree(node *SubNodeType) map[string]SubTreeDebugShow {
	childObj := make(map[string]SubTreeDebugShow)
	node.ChildForEach(func(ntype, nvalue string, nd *SubNodeType) {
		key := fmt.Sprintf("%s.%s", ntype, nvalue)
		childObj = c.PrepareSubNode(nd, key, childObj)
	})
	return childObj
}

func (c *CacheManagerSubTreeMap) PrepareSubNode(node *SubNodeType, ntype string, output map[string]SubTreeDebugShow) map[string]SubTreeDebugShow {
	output[ntype] = SubTreeDebugShow{
		node,
		node.GetCbfnNodeLength(),
		node.CbfnLinkLength(),
		c.setSubtree(node),
	}

	return output
}

func (c *CacheManagerSubTreeMap) prepareSubtreeCache() map[string]SubTreeDebugShow {
	output := make(map[string]SubTreeDebugShow)
	c.data.Range(func(k, v interface{}) bool {
		nodeType := k.(string)
		s := v.(*SubNodeType)
		output = c.PrepareSubNode(s, nodeType, output)
		return true
	})

	return output
}

func (cm *CacheManager) DumpCachedNodes() map[string]interface{} {
	m := make(map[string]interface{})
	m["subtree"] = cm.subTree.prepareSubtreeCache()

	subList := map[string]uint32{}
	cm.subList.data.Range(func(k, v interface{}) bool {
		subVal := v.(uint32)
		subList[fmt.Sprint(k)] = subVal
		return true
	})
	m["subList"] = subList

	nodeCache := map[string]*base.BaseNode{}
	cm.nodeCache.data.Range(func(key, value interface{}) bool {
		v := value.(*base.BaseNode)
		nodeCache[fmt.Sprint(key)] = v
		return true
	})

	m["nodecache"] = nodeCache

	return m
}
