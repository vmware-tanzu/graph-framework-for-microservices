FROM gcr.io/nsx-sm/tools:latest

RUN curl -LJO https://get.helm.sh/helm-v3.9.1-linux-amd64.tar.gz && \
    tar -xvf helm-v3.9.1-linux-amd64.tar.gz && \
    cp linux-amd64/helm /usr/bin/helm && \
    rm helm-v3.9.1-linux-amd64.tar.gz && \
    rm -rf linux-amd64

ARG CHART_NAME

COPY ${CHART_NAME}.tgz /chart.tgz