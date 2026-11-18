# The same fixed official image as tools/bashy-build/image.lock.
ARG BUILD_IMAGE=localhost/bashy-build:v1.0.0-rc.37-go1.27.2-c99fcb7318b2
FROM ${BUILD_IMAGE}
USER 0:0
COPY . /work
RUN chown -R 10001:10001 /work
USER 10001:10001
WORKDIR /work
ENV GOWORK=off GOTOOLCHAIN=local GOMAXPROCS=4 GOFLAGS=-p=1 GOMEMLIMIT=2GiB
RUN /bashy dag --no-journal --force --file DAG.md build-bash
ENTRYPOINT ["/bashy"]
CMD ["dag", "--no-journal", "--force", "--file", "DAG.md", "test-go"]
