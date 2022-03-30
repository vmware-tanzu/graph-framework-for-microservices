ARG BUILDER_TAG

FROM nexus-validation-builder:$BUILDER_TAG
ADD cmd/nexus-validation/nexus-validation /usr/local/bin
CMD nexus-validation