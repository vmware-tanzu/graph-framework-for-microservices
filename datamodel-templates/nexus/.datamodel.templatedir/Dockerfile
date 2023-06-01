FROM gcr.io/nsx-sm/photon-kubectl:latest
ARG IMAGE_NAME
ARG NAME
COPY build/crds /crds
COPY build/nexus-gql/server /build/server
COPY build/install-validator/install.so /build/install.so
RUN echo ${IMAGE_NAME} > /IMAGE
RUN echo ${NAME} > /NAME