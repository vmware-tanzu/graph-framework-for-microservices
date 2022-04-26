FROM harbor-repo.vmware.com/nexus/golang:1.17
ENV GOPRIVATE gitlab.eng.vmware.com
RUN mkdir -p /root/.ssh && \
    chmod -R 0600 /root/.ssh && \
    tdnf install --refresh -y openssh make awk