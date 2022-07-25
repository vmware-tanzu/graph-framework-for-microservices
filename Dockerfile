FROM gcr.io/nsx-sm/tools:latest
ARG IMAGE_NAME
ARG NAME
COPY build/crds /crds
RUN echo ${IMAGE_NAME} > /IMAGE
RUN echo ${NAME} > /NAME