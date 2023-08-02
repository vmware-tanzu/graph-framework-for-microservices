# Upgrading the Datamodel

## Purpose

The purpose of Datamodel Upgrade is to make it possible to upgrade Datamodel without having to first uninstall and then reinstall it.

## What is backward compatibility check?
The Backward Compatibility Check is a validation process performed during a datamodel upgrade to ensure that the changes made to the datamodel is backward compatible with the previously released/installed datamodel.
<br> The check is performed at build time and install time


### Build
#### Parameters

| Name               | Description                                   |  Type  | Value     |
|--------------------|-----------------------------------------------|:------:|:----------|
| `artifact_repo`    | The git repo url where all current CRDs exist | string | `git_url` |
| `prev-spec-branch` | Branch of current CRDs to compare to new CRDs | string | `master`  |

* If "prev-spec-branch" is provided and artifact_repo is provided, datamodel changes are compared against the datamodel artifacts present in prev_spec_branch of artifact_repo provided
*  If the prev-spec-branch is provided but the artifact_repo is not provided, the nexus.yaml file is checked for the artifact_repo information, the built datamodel will be compared with datamodel artifacts present in prev_spec_branch of artifact_repo provided.
*  If the prev-spec-branch is provided but the artifact_repo is not available, an error is thrown and the backward compatibility check is flagged as failed.
<br></br>  
#### When does backward compatibility check fail?
When the change made to the DSL is determined as backward-incompatible (as defined below) then the backward compatibility check in CI will fail.<br>
<br> To enable backward compatibility check at build time we should provide prev-spec-branch and artifact_repo info to build command as shown below<br>
```
$ nexus datamodel build --name orgchart  --artifact_repo=<git_repo> --prev_spec_branch=<master>
```
#### What is considered a non-breaking change?
Examples of non-breaking change
#####  1. Add an optional field in the nexus node (field with "omitempty" tag )

<details>
<summary>Show example</summary>

1. Add an optional field called `Location` (field with `omitempty` tag) in the node
    ```
    type Leader struct {
        // Tags "Root" as a node in datamodel graph
        nexus.Node

       Name          string
       Designation   string
       DirectReports Manager `nexus:"children"`
    +  Location      string  `json:"location,omitempty"`
    }

2. Rebuild your datamodel
   ```
   nexus datamodel build --name orgchart --prev_spec_branch=master
   ```

   Now, the build would succeed

</details>

[//]: # (#####  2. Modify the REST APIs URIs)

[//]: # ()
[//]: # (<details>)

[//]: # (<summary>Show example</summary>)

[//]: # ()
[//]: # (1. Remove the GET URI from the spec)

[//]: # (    ```shell)

[//]: # (   var LeaderRestAPISpec = nexus.RestAPISpec{)

[//]: # (    Uris: []nexus.RestURIs{)

[//]: # (    -   {)

[//]: # (    -       Uri:     "/leader/{root.Leader}",)

[//]: # (    -       Methods: nexus.DefaultHTTPMethodsResponses,)

[//]: # (    -   },)

[//]: # (        {)

[//]: # (            Uri:     "/leaders",)

[//]: # (            Methods: nexus.HTTPListResponse,)

[//]: # (        },)

[//]: # (    },)

[//]: # (   })

[//]: # ()
[//]: # (   // nexus-rest-api-gen:LeaderRestAPISpec)

[//]: # (   type Leader struct {)

[//]: # (      nexus.Node    Name          string)

[//]: # (      Designation   string)

[//]: # (      DirectReports Manager `nexus:"children"`)

[//]: # (      Location      string  `json:"location,omitempty"`)

[//]: # (   })

[//]: # ()
[//]: # (2. Rebuild your datamodel)

[//]: # (   ```)

[//]: # (   nexus datamodel build --name orgchart --prev_spec_branch master)

[//]: # (   ```)

[//]: # ()
[//]: # (   Now, the build would succeed)

[//]: # ()
[//]: # (</details>)

[//]: # ()
[//]: # (#####  3. Modify the REST APIs annotation)

[//]: # ()
[//]: # (<details>)

[//]: # (<summary>Show example</summary>)

[//]: # ()
[//]: # (1. Rename the `nexus-rest-api-gen` annotation spec from `LeaderRestAPISpec` to `NewLeaderRestAPISpec`)

[//]: # ()
[//]: # (   ```shell)

[//]: # (   var NewLeaderRestAPISpec = nexus.RestAPISpec{)

[//]: # (    Uris: []nexus.RestURIs{)

[//]: # (        {)

[//]: # (            Uri:     "/leaders",)

[//]: # (            Methods: nexus.HTTPListResponse,)

[//]: # (        },)

[//]: # (    },)

[//]: # (   })

[//]: # ()
[//]: # (   // nexus-rest-api-gen:NewLeaderRestAPISpec <==)

[//]: # (   type Leader struct {    )

[//]: # (      nexus.Node    Name          string)

[//]: # (      Designation   string)

[//]: # (      DirectReports Manager `nexus:"children"`)

[//]: # (      Location      string  `json:"location,omitempty"`)

[//]: # (   })

[//]: # ()
[//]: # (2. Rebuild your datamodel)

[//]: # (   ```)

[//]: # (   nexus datamodel build --name orgchart --prev_spec_branch master)

[//]: # (   ```)

[//]: # ()
[//]: # (   Now, the build would succeed)

[//]: # ()
[//]: # (</details>)

[//]: # (<br>)

##### What is considered a backward incompatible change?

#####  1. Add a required field in the nexus node

<details>
<summary>Show example</summary>

1. Add a required field called `AdditionalField` in the existing nexus node
   ```
    type Leader struct {
        // Tags "Root" as a node in datamodel graph
        nexus.Node

        Name            string
        Designation     int
        DirectReports   Manager `nexus:"children"`
        Location        string  `json:"location,omitempty"`
    +   AdditionalField string
    }

2. Rebuild the datamodel
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

3. Use the `—force=true` flag to ignore any build failures and obtain successful code generation.
   ```
   nexus datamodel build --name orgchart --force=true
   ```

</details>

#####  2. Remove a required field from the nexus node

<details>
<summary>Show example</summary>

1. Remove a required field called `AdditionalField` from the `Leader` node
   ```
    type Leader struct {
        // Tags "Root" as a node in datamodel graph
        nexus.Node

        Name            string
        Designation     int
        DirectReports   Manager `nexus:"children"`
        Location        string  `json:"location,omitempty"`
    -   AdditionalField string
    }

2. Rebuild the datamodel
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
   time="2023-01-24T14:08:00+05:30" level=error msg="\t >     additionalField:\n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t >       type: string\n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t >     \n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t >   \n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t > \n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t > /spec/versions/name=v1/schema/openAPIV3Schema/properties/spec/required\n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t >   - one required field removed:\n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t >     - additionalField\n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t >     \n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t >   \n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t > \n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t > \n"
   time="2023-01-24T14:08:00+05:30" level=error msg="\t > \n"

3. Use the `—force=true` flag to ignore any build failures and obtain successful code generation.
   ```
   nexus datamodel build --name orgchart --force=true
   ```

</details>

#####  3. Remove an optional field from the nexus node (field with "omitempty" tag )

<details>
<summary>Show example</summary>

1. Remove an optional field called `Location` from the `Leader` node
   ```
    type Leader struct {
        // Tags "Root" as a node in datamodel graph
        nexus.Node

        Name            string
        Designation     int
        DirectReports   Manager `nexus:"children"`
    -   Location        string  `json:"location,omitempty"`
    }

2. Rebuild the datamodel
   ```
   nexus datamodel build --name orgchart 
   ```

   Now, the build would fail and display the incompatible changes as shown below.

   ```
   panic: Error occurred when checking datamodel compatibility: datamodel upgrade failed due to incompatible datamodel changes: \n"
   time="2023-01-24T20:56:26+05:30" level=error msg="\t >  detected changes in model stored in leaders.root.orgchart.org\n"
   time="2023-01-24T20:56:26+05:30" level=error msg="\t > \n"
   time="2023-01-24T20:56:26+05:30" level=error msg="\t > spec changes: \n"
   time="2023-01-24T20:56:26+05:30" level=error msg="\t > /spec/versions/name=v1/schema/openAPIV3Schema/properties/spec/properties\n"
   time="2023-01-24T20:56:26+05:30" level=error msg="\t >   - one field removed:\n"
   time="2023-01-24T20:56:26+05:30" level=error msg="\t >     location:\n"
   time="2023-01-24T20:56:26+05:30" level=error msg="\t >       type: string\n"
   time="2023-01-24T20:56:26+05:30" level=error msg="\t >     \n"
   time="2023-01-24T20:56:26+05:30" level=error msg="\t >   \n"
   time="2023-01-24T20:56:26+05:30" level=error msg="\t > \n"
   time="2023-01-24T20:56:26+05:30" level=error msg="\t > \n"
   time="2023-01-24T20:56:26+05:30" level=error msg="\t > \n"

3. Use the `—force=true` flag to ignore any build failures and obtain successful code generation.
   ```
   nexus datamodel build --name orgchart --force=true
   ```

</details>

<br>

#### How to perform force upgrade ?
If "prev-spec-branch" is not provided, backward compatibility check will be skipped.

### Install

#### Parameters

| Name                                         | Description                                         | Type | Value   |
|----------------------------------------------|-----------------------------------------------------|:----:|:--------|
| `force`                                      | Force upgrade of a datamodel                        | bool | `false` |
| `datamodel_backward_compatibility_validator` | Flag to enable/disable backward compatibility check | bool | `false` |


<br>At runtime:<br>
* Datamodel objects existing in the cluster will be the source of truth 



## How to perform force upgrade at runtime?
1. Ensure there are no datamodel objects for the datamodel node crd that is being upgraded





