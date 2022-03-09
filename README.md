# nexus-sdk

This README and whole nexus-sdk is still in very early point, work in progress...

 nexus-sdk main responsibilities are:
- ...
- ...

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

2. Build nexus-sdk: `make build_in_container`

### Build in custom/local environment

To build nexus-sdk on custom/local environment: `make build`

## Test
### Run tests in containerized sandbox (Recommended)

To run tests in a fixed/sandboxed environment:

1. Download the test sandbox: `make docker.builder`

2. Test nexus-sdk: `make test_in_container`

### Run tests in custom/local environment

To test nexus-sdk on custom/local environment:

1. Download the required tools: `make tools`

2. Run tests: `make test`

# Packaging

nexus-sdk is packaged and published a Docker container images.

Packaging is achieved by the following two steps:
## Creating a base image

Create base image: `make docker.base`

The base image contains all the requirements needed for the nexus-sdk to
run. Also, it can be extended with some handy debug tools like `tcpdump`. To
activate this just set `DEBUG=TRUE` in your environment.
## nexus-sdk image

To build nexus-sdk docker image: `make docker`

# Publishing

nexus-sdk docker image can be published by invoking: `make publish`
