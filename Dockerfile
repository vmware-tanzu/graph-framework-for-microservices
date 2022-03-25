# Build the <app-name> binary
FROM gcr.io/nsx-sm/photon:3.0
ARG APP_NAME
WORKDIR /bin
ENV ENV_APP_NAME=${APP_NAME}
COPY /bin/${APP_NAME} /bin
ENTRYPOINT "/bin/$ENV_APP_NAME"