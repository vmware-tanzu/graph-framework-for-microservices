# nexus compiler

This README and whole nexus-compiler is still in very early point, work in progress...

Nexus compiler main responsibility is to generate code based on provided datamodel. Currently, generated are:
- crd yamls
- crd go client
- crd go apis

# Run compiler in container
1. Create basic application structure, your application should be in GOPATH
2. Add datamodel to your application, example structure:
```
.
├── go.mod
├── main.go
└── nexus
    └── datamodel
        ├── go.mod
        ├── config
        │   └── config.go
        ├── inventory
        │   └── inventory.go
        ├── nexus
        │   └── nexus.go
        ├── root.go
        └── runtime
            └── runtime.go
```
3. Download compiler image or build in compiler repo using `make docker.builder && make docker` command
4. Run compiler from your application nexus directory. Specify GROUP_NAME env variable with your CRD group name:
Your datamodel should be mounter to /go/src/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/datamodel directory, and directory to
which you would like to generate your files should be mounted to /go/src/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/generated
directory. CRD_MODULE_PATH env var will determine import paths for genereted files.
If you follow structure from example above you just need to specify GROUP_NAME and copy rest of following example
command
```
$ cd nexus
$ GROUP_NAME=helloworld.com && \
docker run  \
 --volume $(realpath .)/datamodel:/go/src/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/datamodel  \
 --volume $(realpath .)/generated:/go/src/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/generated \
 -e CRD_MODULE_PATH=$(go list -m)/nexus/generated/ \
 -e GROUP_NAME=$GROUP_NAME  \
 --workdir /go/src/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/ \
  nexus-compiler:1ce29d44
```


# Development
## Guidelines

Any and every code change MUST follow code style guidelines and pass following checks:

- Unit Test cases are mandatory for any and every code change.

- `make test-fmt` - makes sure code formatting matches standard golang formatting

- `make lint` and `make vet` - static code analysis looking for possible programming errors, bugs, stylistic errors, and suspicious constructs

- `make race-unit-test` - executes unit tests with race flag to look for possible race conditions

## Build
### Build in containerized sandbox (Recommended)

To run build in a fixed/sandboxed environment:

1. Download the build sandbox: `make docker.builder`

2. Build nexus compiler: `make build_in_container`

### Build in custom/local environment

To build nexus compiler on custom/local environment: `make build`

### Generate code for example datamodel in custom/local environment

Install required tools using `make tools`

To render templates for example datamodel use `make render_templates`.
To generate all code for example use `make generate_example`.

## Test
### Run tests in containerized sandbox (Recommended)

To run tests in a fixed/sandboxed environment:

1. Download the test sandbox: `make docker.builder`

2. Test nexus compiler: `make test_in_container`

3. Test generation with `make test_generate_code_in_container`

### Test CRD templates rendering:

To render crd templates you can run:
`make render_templates`

This will generate rendered templates to `example/output/_crd_base` directory. This directory can be used for unit tests.

# Packaging

nexus compiler is packaged and published a Docker container images.

Packaging is achieved by the following two steps:
## Creating a base image

Create base image: `make docker`


## nexus compiler image

To build nexus compiler docker image: `make docker`

# Publishing

nexus compiler docker image can be published by invoking: `make publish`
