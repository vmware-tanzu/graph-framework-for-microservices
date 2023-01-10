# Getting Started

Start by creating your own extensible, distributed platform that:

* Implements a datamodel as K8s CRDs.

* Generates accessors that faciliate easy consumption of the datamodel.

* Bootstrap a cloud native application that consumes that datamodel and is ready-to-go in no time.

## **Before getting started**

1. Install Nexus CLI
    ```
    curl -fsSL https://raw.githubusercontent.com/vmware-tanzu/graph-framework-for-microservices/main/cli/get-nexus-cli.sh -o get-nexus-cli.sh
    bash get-nexus-cli.sh
    ```
    <details><summary>FAQs</summary>
      
    To install the specific version
    ```
    bash get-nexus-cli.sh --version <version-tag> 
    ``` 
    
    To install the specific version and the user given destination directory
    ```
    bash get-nexus-cli.sh --version <version-tag> --dst_dir <destination-directoy-path>
    ``` 
	
    </details>
2. Verify nexus sdk pre-requisites are satisfied

    nexus prereq verify


## **Ways to get started**

- **[Playground -- Helloworld Datamodel](Helloworld.md)**

- **[App Local Datamodel](AppLocalDatamodel.md)**

- **[Importable Datamodel](ImportableDatamodel.md)**

