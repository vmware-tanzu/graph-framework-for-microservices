package utils

import (
	"encoding/json"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/ifc"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/logging"

	"github.com/gopherjs/gopherjs/js"
)

const ISOTimeFormat = "2006-01-02T15:04:05.999Z"

func ArrayToJ(dt []string) *js.Object {
	r := js.Global.Get("Array").New()
	for idx, v := range dt {
		r.SetIndex(idx, v)
	}
	return r
}
func ArrayToG(dt *js.Object) []string {
	ret := make([]string, dt.Length())
	for idx := 0; idx < dt.Length(); idx++ {
		ret[idx] = dt.Index(idx).String()
	}
	return ret
}
func NodePathToG(npl *js.Object) ifc.NodePathList {
	ret := ifc.NodePathList{}
	for idx := 0; idx < npl.Length(); idx++ {
		itm := ifc.NodePath{"", ""}
		itmIn := npl.Index(idx)
		itm[ifc.NodePathName_nodeIdentifier] = itmIn.Index(ifc.NodePathName_nodeIdentifier).String()
		itm[ifc.NodePathName_nodeType] = itmIn.Index(ifc.NodePathName_nodeType).String()
		ret = append(ret, itm)
	}
	return ret
}
func NodePathToJ(npl ifc.NodePathList) *js.Object {
	ret := js.Global.Get("Object").New()
	for idx, key := range npl {
		dt := js.Global.Get("Object").New()
		dt.SetIndex(ifc.NodePathName_nodeIdentifier, key[ifc.NodePathName_nodeIdentifier])
		dt.SetIndex(ifc.NodePathName_nodeType, key[ifc.NodePathName_nodeType])
		ret.SetIndex(idx, dt)
	}
	return ret
}

func PropToJ(prop ifc.PropertyType) *js.Object {
	p := js.Global.Get("Object").New()
	for k, v := range prop {
		p.Set(k, v)
	}
	return p
}
func PropToG(prop *js.Object) ifc.PropertyType {
	p := make(ifc.PropertyType)
	for _, key := range js.Keys(prop) {
		val := prop.Get(key)
		p[key] = val.Interface()
	}
	return p
}
func ArrayFromStr(str string) []string {
	ret := []string{}
	if err := json.Unmarshal([]byte(str), &ret); err != nil {
		panic(err)
	}
	return ret
}

func JsonMarshal(dt interface{}) []byte {
	b, err := json.Marshal(dt)
	if err != nil {
		logging.Fatalf("Error when runningjson.Marshal(%s). Error = %s", dt, err)
		panic(err)
	}
	return b
}

func ArrayToStr(dt []string) string {
	b, err := json.Marshal(dt)
	if err != nil {
		panic(err)
	}
	return string(b)
}
func PropFromStr(str string) ifc.PropertyType {
	var prop ifc.PropertyType = ifc.PropertyType{}
	if err := json.Unmarshal([]byte(str), &prop); err != nil {
		panic(err)
	}
	return prop
}
func PropToStr(prop ifc.PropertyType) string {
	b, err := json.Marshal(prop)
	if err != nil {
		panic(err)
	}
	return string(b)
}
func NPathFromStr(str string) ifc.NodePathList {
	prop := ifc.NodePathList{}
	if err := json.Unmarshal([]byte(str), &prop); err != nil {
		panic(err)
	}
	return prop
}
func NPathToStr(prop ifc.NodePathList) string {
	b, err := json.Marshal(prop)
	if err != nil {
		panic(err)
	}
	return string(b)
}
func GLinkToStr(prop *ifc.GLink) string {
	b, err := json.Marshal(prop)
	if err != nil {
		panic(err)
	}
	return string(b)
}
func GLinkToG(lnk *js.Object) *ifc.GLink {
	r := &ifc.GLink{
		Id:                lnk.Get("id").String(),
		LinkType:          lnk.Get("linkType").String(),
		Properties:        PropToG(lnk.Get("properties")),
		SourceNodeId:      lnk.Get("sourceNodeId").String(),
		DestinationNodeId: lnk.Get("destinationNodeId").String()}
	return r
}
func GLinkToJ(lnk *ifc.GLink) *js.Object {
	gl := js.Global.Get("Object").New()
	gl.Set("id", lnk.Id)
	gl.Set("linkType", lnk.LinkType)
	gl.Set("properties", lnk.Properties)
	gl.Set("sourceNodeId", lnk.SourceNodeId)
	gl.Set("destinationNodeId", lnk.DestinationNodeId)
	return gl
}
func GNodeToG(dt *js.Object) *ifc.GNode {
	r := &ifc.GNode{
		Id:         dt.Get("id").String(),
		Type:       dt.Get("type").String(),
		Properties: PropToG(dt.Get("properties"))}
	jlinks := dt.Get("links")
	for i := 0; i < jlinks.Length(); i++ {
		r.Links = append(r.Links, GLinkToG(jlinks.Index(i)))
	}
	return r
}
