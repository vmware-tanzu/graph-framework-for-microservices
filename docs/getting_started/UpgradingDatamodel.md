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

## What is considered a non-breaking change?

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
   nexus datamodel build --name orgchart
   ```

   Now, the build would succeed

</details>

#####  2. Modify the REST APIs URIs

<details>
<summary>Show example</summary>

1. Remove the GET URI from the spec
    ```shell
   var LeaderRestAPISpec = nexus.RestAPISpec{
    Uris: []nexus.RestURIs{
    -   {
    -       Uri:     "/leader/{root.Leader}",
    -       Methods: nexus.DefaultHTTPMethodsResponses,
    -   },
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
      Location      string  `json:"location,omitempty"`
   }

2. Rebuild your datamodel
   ```
   nexus datamodel build --name orgchart
   ```

   Now, the build would succeed

</details>

#####  3. Modify the REST APIs annotation

<details>
<summary>Show example</summary>

1. Modify the `nexus-rest-api-gen` annotation spec from `LeaderRestAPISpec` to `NewLeaderRestAPISpec`
    
   ```shell
   var NewLeaderRestAPISpec = nexus.RestAPISpec{
    Uris: []nexus.RestURIs{
        {
            Uri:     "/leaders",
            Methods: nexus.HTTPListResponse,
        },
    },
   }

   // nexus-rest-api-gen:NewLeaderRestAPISpec <==
   type Leader struct {    
      nexus.Node    Name          string
      Designation   string
      DirectReports Manager `nexus:"children"`
      Location      string  `json:"location,omitempty"`
   }

2. Rebuild your datamodel
   ```
   nexus datamodel build --name orgchart 
   ```

   Now, the build would succeed

</details>

## What is considered a backward incompatible change?

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
