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

The following cases are considered non-breaking and backward incompatible changes:

| Use Cases                                             | force=false | force=true |
|-------------------------------------------------------|:-----------:|:----------:|
| Add/Remove a field in the nexus node                  |   &cross;   |  &check;   |
| Add a optional field ("omitempty" tag )               |   &check;   |  &check;   |
| Remove a optional field ("omitempty" tag )            |   &cross;   |  &check;   |
| Add/Remove a child/link                               |   &cross;   |  &check;   |
| Modify the type of a field                            |   &cross;   |  &check;   |
| Modify the field name                                 |   &cross;   |  &check;   |
| Modify the REST API URIs                              |   &check;   |  &check;   |
| Add a field in status sub-resource of a nexus node    |   &check;   |  &check;   |
| Remove a field in status sub-resource of a nexus node |   &cross;   |  &check;   |
| Modify the field in status sub-resource               |   &cross;   |  &check;   |


* To achieve a successful installation with force=true, you must manually delete any CR objects that are already in the system.


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
