# integration-tests

### Cleanup buckets with the commit ID created from CI

To list the folders and check if they are created before 10 days
```
python cleanup_bucket.py nexus-template-downloads 10
```

To delete the listed folders.

```
python cleanup_bucket.py nexus-template-downloads 10 delete
```

### Running integration tests in CI for new component.

Step 1:


Add the include at top of .gitlab-ci.yml of the repo
```
include:
- project: "nsx-allspark_users/nexus-sdk/integration-tests"
  ref: NPT-171
  file: ".gitlab-ci-template.yml"
```

Step 2:

Append this at bottom of .gitlab-ci.yml file
```
integration_test:
  stage: integration_test
  variables:
    INTEGRATION_TEST_BRANCH: NPT-171
    NEXUS_<COMPONENT>_VERSION: $CI_COMMIT_SHA
  only:
    - merge_requests
    - master
  extends: .run_integration_test
```
Add the below line to stages in .gitlab-ci.yml file
```
- integration_test
```

* Note: NEXUS_<COMPONENT>_VERSION this needs to be added in cli repo code

Example to add runtime manifests
```
runtimeVersion := os.Getenv("NEXUS_RUNTIME_TEMPLATE_VERSION")
if runtimeVersion == "" {
		runtimeVersion = values.NexusDatamodelTemplates.Version
}
```

Optional:

To test with Custom CLI branch
```
integration_test:
  stage: integration_test
  variables:
    INTEGRATION_TEST_BRANCH: NPT-171
    NEXUS_<COMPONENT>_VERSION: $CI_COMMIT_SHA
    CLI_BRANCH: <<custom_branch>>
  only:
    - merge_requests
    - master
  extends: .run_integration_test
```
Update integration test variables , revert back before merging.

To debug integration-tests failures:
Step 1:

Add --debug for the command(s) on docs repo

Step 2:

Add the DOCS_BRANCH variable with the custom docs branch.

```
integration_test:
  stage: integration_test
  variables:
    INTEGRATION_TEST_BRANCH: NPT-171
    NEXUS_<COMPONENT>_VERSION: $CI_COMMIT_SHA
    CLI_BRANCH: <<custom_branch>>
    DOCS_BRANCH: <<custom_branch_doc>>   <----- add the doc branch where you added the --debug lines
  only:
    - merge_requests
    - master
  extends: .run_integration_test
```


## To run the integration tests locally

Pre-Requisites:

Access to gitlab nexus-sdk group.
#### Step0(Optional): For changes present in private branch for microservices

### For image update

Update Nexus repository with the commit SHA(s) of the microservice(s) in this section:

`https://gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/nexus/-/blob/master/nexus-runtime/values.yaml#L7-16`

### For changes present in helm charts of any microservice or runtime manifests

Update Nexus repository submodule definition

`https://gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/nexus/-/blob/master/.gitmodules#L1-20`

Edit branch section and run

```
make submodule
```

Commit to a private branch in Nexus repository


`Note: Please create draft MR in microservice(s) and Nexus repo to build images and helm charts and push it.`

### Update CLI Repository

#### For Changes in Nexus repository

Please update https://gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli/-/blob/master/pkg/common/values.yaml#L10
with v0.0.0-<commit-SHA> if changes are pushed to private branch and Draft MR created.

Please update https://gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli/-/blob/master/pkg/common/values.yaml file with commit ID based on component you want to test.

Commit and push to remote private branch

#### Step1: Create kind cluster to install runtime

Please create a kind k8s cluster to start integration tests
```
kind create cluster
```

####  Step2: Install CLI from the branch you need to test

for running integration tests from master CLI
```
export CLI_BRANCH=master
```

for custom branch replace branch name in above export
```
export CLI_BRANCH=<branch>
```

Please check if there are folders with existing files

```
$GOPATH/test-app
$GOPATH/test-app-local
$GOPATH/test-app-imported-dm
```

If any of above directories present , please take a backup and copy it to different folder or remove this folders

```
rm -rf $GOPATH/test-app $GOPATH/test-app-local $GOPATH/test-app-imported-dm
```

#### Step3: Check Pre-requisite environment variables present
```
export GOPATH=<>
export KIND_CLUSTER_NAME=kind
export GOPRIVATE=*.eng.vmware.com
export GOINSECURE=*.eng.vmware.com
```

`Note: Please use the kind cluster name you created above`

```
GOPRIVATE="gitlab.eng.vmware.com" go install gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/cmd/plugin/nexus@$CLI_BRANCH
```

#### Step4: Start Integration test
```
make run_integration_test
```