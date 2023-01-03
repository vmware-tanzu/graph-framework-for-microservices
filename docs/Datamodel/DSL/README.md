Nexus DSL provides the framework and syntax for declaring Nexus datamodel nodes and types using Golang syntax.

```
/*  Section 1: Group Declaration */

package gns                                                                       <--- API / Node group name


/* Section 2: Attribute Definition (OPTIONAL) */

var GNSRestAPISpec = nexus.RestAPISpec{
	Uris: []nexus.RestURIs{                                                       <--- List of REST URL on which the Nexus Node should be exposed
		{
			Uri:     "/v1alpha2/projects/{project}/global-namespace/{gns.Gns}",   <--- REST URL on which the Nexus Node should be exposed
			Methods: nexus.HTTPMethodsResponses{                                  <--- Methods and responses to be enabled on this REST URL
				http.MethodGet: nexus.DefaultHTTPGETResponses,
			},
		},
		{
			Uri:     "/v1alpha2/global-namespace",                                <--- REST URL on which the Nexus Node should be exposed
			QueryParams: []string{
				"project",                                                        <--- Instead of URI param we are using QueryParams to specify project
			    "gns.Gns"
			},
			Methods: nexus.HTTPMethodsResponses{                                  <--- Methods and responses to be enabled on this REST URL
				http.MethodGet: nexus.DefaultHTTPGETResponses,
			},
		},
		{
			Uri:     "/v1alpha2/projects/{project}/global-namespaces",            <--- REST URL on which the Nexus Node should be exposed
			Methods: nexus.HTTPListResponse,                                      <--- nexus.HTTPListResponse indicates that this request will return a list of objects
		},
	},
}

/* Section 3: Attribute Binding (OPTIONAL) */
// nexus-rest-api-gen:GNSRestAPISpec                            <--- Binds a REST API attribute with a Nexus Node
// nexus-description: This is a my GNS node description         <--- Adds custom description to Nexus Node. This custom description will be propagated to references to this node in OpenAPI spec.

/* Section 4: Node Definition */
type Gns struct {

    nexus.Node                                                  <--- Declares type "Gns" to be a Nexus Node
                                                                     Alternatively nexus.SingletonNode can be used. SingletonNodes can only have 'default' as a display name
                                                                     For root level nodes there can be only one singleton node present in the system, for non-root objects
                                                                     only one can be present for given parents.


    F1 string                                                   <--- Defines a field of standard type.
                                                                     Supported standard types: bool, int32, float32, string

    F2 string `json:"f2,omitempty"`                             <--- Defines an optional field.
                                                                     To make a field optional, omitempty tag should be added.

    F3 CustomType                                               <--- Defines a field of custom type. The type definition can in the same go package or
                                                                     can be imported from other packages. It should be resolvable by Go compiler.   

    C1 ChildType1 `nexus:"child"`                               <--- Declares:
                                                                     * a child of type "ChildType1"
                                                                     * field C1 through which a specific instance/object of type "ChildType1" can be accessed

                                                                     ChildType1 should be another nexus node in the graph.
                                                                     ChildType1 should be resolvable by Go compiler either in the local package or through Go import

                                                                     C1 is the handle through which a specific object of type ChildType1 can be accessed.
                                                                     While there can be multiple objects of type ChildType1 in the system, C1 can only hold single object.

    C2 ChildType2 `nexus:"children"`                            <--- Declares:
                                                                     * a child of type "ChildType2"
                                                                     * field C2 through which multiple objects/instances of type "ChildType2" can be accessed, with each object queryable by name

                                                                     ChildType2 should be another nexus node in the graph.
                                                                     ChildType2 should be resolvable by Go compiler either in the local package or through Go import.

                                                                     C2 is the handle through which multiple objects of type ChildType2 can be accessed.
                                                                     Objects in C2 are queryable by name.
    
    L1 LinkType1 `nexus:"link"`                                 <--- Declares:
                                                                     * a link of type "LinkType1"
                                                                     * field L1 through which a specific instance/object of type "LinkType1" can be accessed

                                                                     LinkType1 should be another nexus node in the graph.
                                                                     LinkType1 should be resolvable by Go compiler either in the local package or through Go import.

                                                                     L1 is the handle through which a specific object of type LinkType1 can be accessed.
                                                                     While there can be multiple objects of type LinkType1 in the system, L1 can only hold single object.

    L2 LinkType2 `nexus:"links"`                                <--- Declares:
                                                                    * a link of type "LinkType2"
                                                                    * field L2 through which multiple objects/instances of type "LinkType2" can be accessed, with each object queryable by name

                                                                     LinkType2 should be another nexus node in the graph.
                                                                     LinkType2 should be resolvable by Go compiler either in the local package or through Go import.

                                                                     L2 is the handle through which multiple objects of type LinkType2 can be accessed.
                                                                     Objects in L2 are queryable by name.

    S1 StatusType1 `nexus:"status"`                             <--- Declares a status field of type "StatusType".
}
```
