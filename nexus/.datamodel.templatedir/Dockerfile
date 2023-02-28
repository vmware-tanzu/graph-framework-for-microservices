FROM gcr.io/nsx-sm/tools:latest
ARG IMAGE_NAME
ARG NAME
COPY build/crds /crds
COPY build/nexus-gql/graphql.so /build/graphql.so
COPY build/install-validator/install.so /build/install.so
RUN echo ${IMAGE_NAME} > /IMAGE
RUN echo ${NAME} > /NAME