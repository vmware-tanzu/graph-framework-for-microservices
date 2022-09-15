ARG BUILDER_TAG


FROM gcr.io/nsx-sm/photon:4.0
ENV GOPATH /go
WORKDIR /usr/local/bin
# Need to be enabled for supporting load plugins
ENV CGO_ENABLED 1
ADD main /usr/local/bin
ENV GOLANG_PROTOBUF_REGISTRATION_CONFLICT warn
CMD ["./main"]



