package model

import (
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/common-library.git/pkg/nexus"
	"golang.org/x/net/publicsuffix"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sync"
)

var (
	RestURIChan = make(chan []nexus.RestURIs, 100)
	CrdTypeChan = make(chan string, 100)

	// OidcChan is used to pass on OIDC node updates to the OIDC authenticator
	OidcChan = make(chan OidcNodeEvent)

	CrdTypeToRestUris      = make(map[string][]nexus.RestURIs)
	crdTypeToRestUrisMutex = &sync.Mutex{}

	// CRD name to CRD type (Gns.gns => gns.vmware.org)
	UriToCRDType      = make(map[string]string)
	uriToCRDTypeMutex = &sync.Mutex{}

	// CRD Type to NodeInfo (gns.vmware.org => NodeInfo{})
	CrdTypeToNodeInfo      = make(map[string]NodeInfo)
	crdTypeToNodeInfoMutex = &sync.Mutex{}

	// CRD Type to k8s spec (gns.vmware.org => CustomResourceDefinitionSpec)
	CrdTypeToSpec      = make(map[string]apiextensionsv1.CustomResourceDefinitionSpec)
	crdTypeToSpecMutex = &sync.Mutex{}

	DatamodelsChan                = make(chan string, 100)
	DatamodelToDatamodelInfo      = make(map[string]DatamodelInfo)
	DatamodelToDatamodelInfoMutex = &sync.Mutex{}
)

func ConstructDatamodel(eventType EventType, name string, unstructuredObj *unstructured.Unstructured) {
	DatamodelToDatamodelInfoMutex.Lock()
	defer DatamodelToDatamodelInfoMutex.Unlock()

	if eventType == Delete {
		delete(DatamodelToDatamodelInfo, name)
		return
	}
	obj := unstructuredObj.Object

	spec := obj["spec"].(map[string]interface{})

	if title, ok := spec["title"]; ok {
		//FIXME: data race
		datamodelName, _ := publicsuffix.EffectiveTLDPlusOne(name)
		DatamodelToDatamodelInfo[datamodelName] = DatamodelInfo{
			Title: title.(string),
		}

		DatamodelsChan <- datamodelName
	}
}

func ConstructMapURIToCRDType(eventType EventType, crdType string, apiURIs []nexus.RestURIs) {
	uriToCRDTypeMutex.Lock()
	defer uriToCRDTypeMutex.Unlock()

	if eventType == Delete {
		for uri, cType := range UriToCRDType {
			if cType == crdType {
				delete(UriToCRDType, uri)
			}
		}
	}

	for _, u := range apiURIs {
		UriToCRDType[u.Uri] = crdType
	}
}

func ConstructMapCRDTypeToNode(eventType EventType, crdType, name string, parentHierarchy []string,
	children map[string]NodeHelperChild, isSingleton bool, description string) {
	crdTypeToNodeInfoMutex.Lock()
	defer crdTypeToNodeInfoMutex.Unlock()

	if eventType == Delete {
		delete(CrdTypeToNodeInfo, crdType)
	}

	CrdTypeToNodeInfo[crdType] = NodeInfo{
		Name:            name,
		ParentHierarchy: parentHierarchy,
		Children:        children,
		IsSingleton:     isSingleton,
		Description:     description,
	}

	// Push new CRD Type to chan
	CrdTypeChan <- crdType
}

func ConstructMapCRDTypeToSpec(eventType EventType, crdType string, spec apiextensionsv1.CustomResourceDefinitionSpec) {
	crdTypeToSpecMutex.Lock()
	defer crdTypeToSpecMutex.Unlock()

	if eventType == Delete {
		delete(CrdTypeToSpec, crdType)
	}
	CrdTypeToSpec[crdType] = spec
}

func ConstructMapCRDTypeToRestUris(eventType EventType, crdType string, restSpec nexus.RestAPISpec) {
	crdTypeToRestUrisMutex.Lock()
	defer crdTypeToRestUrisMutex.Unlock()

	if eventType == Delete {
		delete(CrdTypeToRestUris, crdType)
	}

	uris := CrdTypeToRestUris[crdType]
	uris = append(uris, restSpec.Uris...)
	CrdTypeToRestUris[crdType] = uris

	// Push new uris to chan
	RestURIChan <- uris
}
