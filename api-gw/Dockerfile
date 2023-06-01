# syntax = harbor-repo.vmware.com/dockerhub-proxy-cache/docker/dockerfile:1.3
FROM harbor-repo.vmware.com/nexus/golang:1.19.8 AS builder
ARG GIT_TAG
ARG GIT_HEAD
ARG CICD_TOKEN
ARG APP_NAME
ARG USE_SSH
RUN mkdir -p /root/.ssh && \
   chmod 0700 /root/.ssh
COPY .ssh/ /root/.ssh/
WORKDIR /api-gw
ENV GOPRIVATE *.eng.vmware.com
ENV GOINSECURE *.eng.vmware.com
ENV CGO_ENABLED=0
COPY go.* .
# Required for access to GO artifactory
RUN git config --global credential.helper store
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
#Intialize go deps
RUN \
   if [ "x$USE_SSH" = "xTRUE" ] || [ "x$USE_SSH" = "xtrue" ] ;\
   then \
       tdnf --refresh install -y openssh && \
       chmod -R 0600 /root/.ssh/* && \
       git config --global --add url."git@gitlab.eng.vmware.com:".insteadOf "https://gitlab.eng.vmware.com/" &&\
       go mod download ;\
   else \
       echo "https://gitlab-ci-token:${CICD_TOKEN}@gitlab.eng.vmware.com" >> ~/.git-credentials && \
       git config --global credential.helper store && \
       go mod download ;\
   fi

# Build
RUN --mount=type=cache,target=/root/.cache/go-build \
    GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X gitlab.eng.vmware.com/nsx-allspark_users/lib-go/allspark/health.Sha=$GIT_HEAD -X gitlab.eng.vmware.com/nsx-allspark_users/lib-go/allspark/health.Tag=$GIT_TAG" -o bin/api-gw main.go

FROM gcr.io/nsx-sm/photon:4.0
WORKDIR /bin
ARG APP_NAME
COPY --from=builder /api-gw/bin/api-gw .
USER 65532:65532
ENTRYPOINT "/bin/api-gw"

