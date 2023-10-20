FROM gcr.io/nsx-sm/photon:4.0
WORKDIR /bin
ARG APP_NAME
COPY /api-gw/bin/api-gw .
USER 65532:65532
ENTRYPOINT "/bin/api-gw"
