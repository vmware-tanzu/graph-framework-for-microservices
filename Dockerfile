
# Build the <app-name> binary
FROM 284299419820.dkr.ecr.us-west-2.amazonaws.com/golang:1.17 AS builder
ARG GIT_TAG
ARG GIT_HEAD
ARG CICD_TOKEN
ARG APP_NAME
WORKDIR /api-gw
ENV GOPRIVATE gitlab.eng.vmware.com
# Required for access to GO artifactory
RUN git config --global credential.helper store
RUN echo "https://gitlab-ci-token:${CICD_TOKEN}@gitlab.eng.vmware.com" >> ~/.git-credentials
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
COPY docs docs
# Copy the go source
COPY main.go main.go
COPY controllers/ controllers/
COPY internal/ internal/
COPY pkg/ pkg/
COPY test/ test/
COPY nexus/ nexus/
# Initialize go deps
RUN go mod tidy && go mod download
# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags "-s -w -X gitlab.eng.vmware.com/nsx-allspark_users/lib-go/allspark/health.Sha=$GIT_HEAD -X gitlab.eng.vmware.com/nsx-allspark_users/lib-go/allspark/health.Tag=$GIT_TAG" -o bin/api-gw main.go

FROM gcr.io/nsx-sm/photon:3.0
WORKDIR /bin
ARG APP_NAME
COPY --from=builder /api-gw/bin/api-gw .
USER 65532:65532
ENTRYPOINT "/bin/api-gw"

