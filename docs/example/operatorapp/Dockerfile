
# Build the <app-name> binary
FROM gcr.io/nsx-sm/photon:3.0
WORKDIR /bin
ARG APP_NAME
COPY bin/test-app-local .
USER 65532:65532
ENTRYPOINT "/bin/test-app-local"

