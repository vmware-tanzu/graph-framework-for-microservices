FROM gcr.io/nsx-sm/nexus-datamodel-installer:v0.0.1
ARG IMAGE_NAME
ARG NAME
COPY build/crds /crds
RUN echo ${IMAGE_NAME} > /IMAGE
RUN echo ${NAME} > /NAME