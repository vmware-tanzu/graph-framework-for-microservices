
# Build the <app-name> binary
FROM 284299419820.dkr.ecr.us-west-2.amazonaws.com/golang:1.17 AS builder
ARG GIT_TAG
ARG GIT_HEAD
ARG CICD_TOKEN
ARG APP_NAME
WORKDIR /${APP_NAME}
ENV GOPRIVATE gitlab.eng.vmware.com
RUN git config --global credential.helper store
RUN echo "https://gitlab-ci-token:${CICD_TOKEN}@gitlab.eng.vmware.com" >> ~/.git-credentials
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# Copy the go source
COPY main.go main.go
COPY controllers/ controllers/
COPY nexus/ nexus/
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod tidy && go mod download
# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags "-s -w -X gitlab.eng.vmware.com/nsx-allspark_users/lib-go/allspark/health.Sha=$GIT_HEAD -X gitlab.eng.vmware.com/nsx-allspark_users/lib-go/allspark/health.Tag=$GIT_TAG" -o bin/${APP_NAME} main.go
# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/nsx-sm/photon:3.0
WORKDIR /bin
ARG APP_NAME
ENV ENV_APP_NAME=${APP_NAME}
COPY --from=builder /${APP_NAME}/bin/${APP_NAME} .
USER 65532:65532
ENTRYPOINT "/bin/$ENV_APP_NAME"

