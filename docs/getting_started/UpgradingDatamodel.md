# Upgrading the Datamodel

## Purpose

The purpose of Datamodel Upgrade is to make it possible to reinstall Datamodel without having to first uninstall and then reinstall it.

## Parameters

### Common parameters(Build + Install)


| Name    | Description                  | Type | Shorthand | Value   |
|---------|------------------------------|:----:|:---------:|:--------|
| `force` | Force upgrade of a datamodel | bool |     f     | `false` |


### Build parameters

| Name               | Description                                           |  Type  | Value                 |
|--------------------|-------------------------------------------------------|:------:|:----------------------|
| `prev-spec-dir`    | The path to the directory containing all current CRDs | string | `build/crds`          |
| `prev-spec-branch` | Branch of current CRDs to compare to new CRDs         | string | `current branch HEAD` |

* if "prev-spec-dir" is provided, (users can provide a directory path), where all the CRD's can be found. We would use those files as a source of truth.
* if "prev-spec-branch" is given, we can use the specified branch to figure out the CRDs directory. Otherwise, default to the HEAD of the current working branch.

Specify each parameter using the `--key=value` argument to `nexus datamodel build`. For example,

```
$ nexus datamodel build --name orgchart --force=false --prev-spec-dir=<build/crds> --prev-spec-branch=<master>
```

The above command builds the new datamodel against the master branch of the build/crds directory.

## What is considered a non-breaking and backward incompatible change?

#### To try the example

1. It requires the Datamodel to be initialised and installed before proceeding to further steps. Please follow the below link to configure datamodel.

   #### [Playground](docs/getting_started/Playground.md)

These are considered non-breaking or backward compatible changes:

#####  1. Add an optional field in the nexus node (field with "omitempty" tag )

<details>
<summary>:eyes: Show example</summary>

1. Edit the root.go file in your datamodel.
2. Add an optional field called `AdditionalField` (field with `omitempty` tag) in the node
    ```shell
   cat <<< '--- root.go.orig	2023-01-05 04:22:07.475996737 -0800
   +++ root.go	2023-01-17 05:41:07.993795597 -0800
   @@ -25,7 +25,8 @@
    type Leader struct {
        // Tags "Root" as a node in datamodel graph
        nexus.Node

   -    Name          string
   -    Designation   string
   -    DirectReports Manager `nexus:"children"`
   +    Name          string
   +    Designation   string
   +    DirectReports Manager `nexus:"children"`
   +    AdditionalField string `json:"additionalField,omitempty"`
    }
   ' | patch root.go --ignore-whitespace
      ```

      <!--
      ```
      cp $DOCS_INTERNAL_DIR/root.go.patched root.go
      ```
      -->

3. Rebuild your datamodel
   ```
   nexus datamodel build --name orgchart
   ```

   Now, the build would succeed

</details>

#####  2. Modify the existing REST API URIs

<details>
<summary>:eyes: Show example</summary>

1. Edit the root.go file in your datamodel.
2. Remove the GET URI from the spec
    ```shell
   var LeaderRestAPISpec = nexus.RestAPISpec{
    Uris: []nexus.RestURIs{
        {
            Uri:     "/leaders",
            Methods: nexus.HTTPListResponse,
        },
    },
   }

   // nexus-rest-api-gen:LeaderRestAPISpec
   type Leader struct {
      nexus.Node    Name          string
      Designation   string
      DirectReports Manager `nexus:"children"`
   }
   ' | patch root.go --ignore-whitespace
      ```

      <!--
      ```
      cp $DOCS_INTERNAL_DIR/root.go.patched root.go
      ```
      -->

3. Rebuild your datamodel
   ```
   nexus datamodel build --name orgchart
   ```

   Now, the build would succeed

</details>

#####  3. Modify the existing REST APIs annotation

<details>
<summary>:eyes: Show example</summary>

1. Edit the root.go file and modify the nexus-rest-api-gen annotation spec 
    
   ```shell
   var NewLeaderRestAPISpec = nexus.RestAPISpec{
    Uris: []nexus.RestURIs{
        {
            Uri:     "/leaders",
            Methods: nexus.HTTPListResponse,
        },
    },
   }

   // nexus-rest-api-gen:NewLeaderRestAPISpec    <====
   type Leader struct {    
      nexus.Node    Name          string
      Designation   string
      DirectReports Manager `nexus:"children"`
   }
   ' | patch root.go --ignore-whitespace
      ```

      <!--
      ```
      cp $DOCS_INTERNAL_DIR/root.go.patched root.go
      ```
      -->

2. Rebuild your datamodel
   ```
   nexus datamodel build --name orgchart 
   ```

   Now, the build would succeed

</details>

And these would be considered breaking changes:

#####  1. Add a new field in the nexus node

<details>
<summary>:eyes: Show example</summary>

1. Edit the root.go file in your datamodel. 
2. Add a new field called `AdditionalField` in the nexus node 
   ```shell
   cat <<< '--- root.go.orig	2023-01-05 04:22:07.475996737 -0800
   +++ root.go	2023-01-17 05:41:07.993795597 -0800
   @@ -25,7 +25,8 @@
    type Leader struct {
        // Tags "Root" as a node in datamodel graph
        nexus.Node

   -    Name          string
   -    Designation   string
   -    DirectReports Manager `nexus:"children"`
   +    Name          string
   +    Designation   int
   +    DirectReports Manager `nexus:"children"`
   +    AdditionalField string
    }
   ' | patch root.go --ignore-whitespace
      ```

      <!--
      ```
      cp $DOCS_INTERNAL_DIR/root.go.patched root.go
      ```
      -->

3. Build your datamodel again
   ```
   nexus datamodel build --name orgchart 
   ```

   Now, the build would fail and display the incompatible changes as shown below.
   ```
   panic: Error occurred when checking datamodel compatibility: datamodel upgrade failed due to incompatible datamodel changes: \n"
   time="2023-01-24T12:21:56+05:30" level=error msg="\t >  detected changes in model stored in leaders.root.orgchart.org\n"
   time="2023-01-24T12:21:56+05:30" level=error msg="\t > \n"
   time="2023-01-24T12:21:56+05:30" level=error msg="\t > spec changes: \n"
   time="2023-01-24T12:21:56+05:30" level=error msg="\t > /spec/versions/name=v1/schema/openAPIV3Schema/properties/spec/required\n"
   time="2023-01-24T12:21:56+05:30" level=error msg="\t >   + one required field added:\n"
   time="2023-01-24T12:21:56+05:30" level=error msg="\t >     - additionalField\n"
   time="2023-01-24T12:21:56+05:30" level=error msg="\t >     \n"
   time="2023-01-24T12:21:56+05:30" level=error msg="\t >   \n"
   time="2023-01-24T12:21:56+05:30" level=error msg="\t > \n"
   time="2023-01-24T12:21:56+05:30" level=error msg="\t > \n"
   time="2023-01-24T12:21:56+05:30" level=error msg="\t > \n"
   
4. Use the `—force=true` flag to ignore any build failures and obtain successful code generation.
   ```
   nexus datamodel build --name orgchart --force=true
   ```
   
</details>

#####  2. Remove a field from the nexus node

<details>
<summary>:eyes: Show example</summary>

1. Edit the root.go file in your datamodel.
2. Remove a field called `Designation` from the `Leader` node  
   ```shell
   cat <<< '--- root.go.orig	2023-01-05 04:22:07.475996737 -0800
   +++ root.go	2023-01-17 05:41:07.993795597 -0800
   @@ -25,7 +25,8 @@
    type Leader struct {
        // Tags "Root" as a node in datamodel graph
        nexus.Node

   -    Name          string
   -    Designation   string
   -    DirectReports Manager `nexus:"children"`
   +    Name          string
   +    DirectReports Manager `nexus:"children"`
    }
   ' | patch root.go --ignore-whitespace
      ```

      <!--
      ```
      cp $DOCS_INTERNAL_DIR/root.go.patched root.go
      ```
      -->

3. Build your datamodel again
   ```
   nexus datamodel build --name orgchart 
   ```

   Now, the build would fail and display the incompatible changes as shown below.

   ```
   panic: Error occurred when checking datamodel compatibility: datamodel upgrade failed due to incompatible datamodel changes: \n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t >  detected changes in model stored in leaders.root.orgchart.org\n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t > \n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t > spec changes: \n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t > /spec/versions/name=v1/schema/openAPIV3Schema/properties/spec/properties\n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t >   - one field removed:\n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t >     designation:\n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t >       type: string\n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t >     \n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t >   \n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t > \n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t > /spec/versions/name=v1/schema/openAPIV3Schema/properties/spec/required\n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t >   - one required field removed:\n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t >     - designation\n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t >     \n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t >   \n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t > \n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t > \n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t > \n"

4. Use the `—force=true` flag to ignore any build failures and obtain successful code generation.
   ```
   nexus datamodel build --name orgchart --force=true
   ```

</details>

#####  3. Remove an optional field from the nexus node (field with "omitempty" tag )

<details>
<summary>:eyes: Show example</summary>

1. Edit the root.go file in your datamodel.
2. Remove an optional field called `Designation` from the `Leader` node
   ```shell
   cat <<< '--- root.go.orig	2023-01-05 04:22:07.475996737 -0800
   +++ root.go	2023-01-17 05:41:07.993795597 -0800
   @@ -25,7 +25,8 @@
    type Leader struct {
        // Tags "Root" as a node in datamodel graph
        nexus.Node

   -    Name          string
   -    Designation   string  `json:"additionalField,omitempty"`
   -    DirectReports Manager `nexus:"children"`
   +    Name          string
   +    DirectReports Manager `nexus:"children"`
    }
   ' | patch root.go --ignore-whitespace
      ```

      <!--
      ```
      cp $DOCS_INTERNAL_DIR/root.go.patched root.go
      ```
      -->

3. Build your datamodel again
   ```
   nexus datamodel build --name orgchart 
   ```

   Now, the build would fail and display the incompatible changes as shown below.

   ```
   panic: Error occurred when checking datamodel compatibility: datamodel upgrade failed due to incompatible datamodel changes: \n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t >  detected changes in model stored in leaders.root.orgchart.org\n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t > \n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t > spec changes: \n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t > /spec/versions/name=v1/schema/openAPIV3Schema/properties/spec/properties\n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t >   - one field removed:\n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t >     designation:\n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t >       type: string\n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t >     \n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t >   \n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t > \n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t > /spec/versions/name=v1/schema/openAPIV3Schema/properties/spec/required\n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t >   - one required field removed:\n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t >     - designation\n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t >     \n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t >   \n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t > \n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t > \n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t > \n"

4. Use the `—force=true` flag to ignore any build failures and obtain successful code generation.
   ```
   nexus datamodel build --name orgchart --force=true
   ```

</details>




## Trying the example

This workflow will walk you through the steps to upgrade your applocal datamodel.

* [Pre-requisites](UpgradingDatamodel.md#pre-requisites)
* [Upgrading](UpgradingDatamodel.md#Upgrading)

### Pre-requisites
1. This workflow requires Datamodel to be initialised and installed before proceeding to further steps. Please follow the below link to configure datamodel

   #### [Playground](docs/getting_started/Playground.md)

### Upgrading

1. Make changes to the root.go file in your datamodel by remove one field, renaming another and adding one.
   These datamodel changes are incompatible (although the newly added field wouldn't be incompatible if it had the `json:"omitempty"` tag )
   ```shell
   cat <<< '--- root.go.orig	2023-01-05 04:22:07.475996737 -0800
   +++ root.go	2023-01-17 05:41:07.993795597 -0800
   @@ -25,7 +25,8 @@
    type Leader struct {
        // Tags "Root" as a node in datamodel graph
        nexus.Node

   -    Name          string
   -    Designation   string
   -    DirectReports Manager `nexus:"children"`
   +    Name          string
   +    Designation   int
   +    DirectReports Manager `nexus:"children"`
   +    AdditionalData string
    }
   ' | patch root.go --ignore-whitespace
      ```

      <!--
      ```
      cp $DOCS_INTERNAL_DIR/root.go.patched root.go
      ```
      -->

2. Build your datamodel again. We have to use `--force=true` flag, or else the build would fail and show us the incompatible changes.
   ```
   nexus datamodel build --name orgchart --force=true
   DOCKER_REPO=orgchart VERSION=latest make docker_build
   kind load docker-image orgchart:latest --name <kind cluster name>
   ```

3. Delete leader object manually, as it is backwards incompatible
   ```shell
   kubectl -s localhost:5000 delete leaders.root.orgchart.org <MyLeader>
   ```

4. Install built datamodel into runtime. Like the build part, we also have to use force flag, or else the upgrade will fail due to incompatibility.
   ```
   nexus datamodel install image orgchart:latest --namespace <name> --force=true
   ```
