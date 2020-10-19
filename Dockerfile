# dynamic config
ARG             BUILD_DATE
ARG             VCS_REF
ARG             VERSION

# build
FROM            golang:1.15-alpine as builder
RUN             apk add --no-cache git gcc musl-dev make bash
ENV             GO111MODULE=on
WORKDIR         /go/src/github.com/MrEhbr/gofsm
COPY            go.* ./
RUN             go mod download
COPY            . ./
RUN             make install

# minimalist runtime
FROM alpine:3.11
LABEL org.label-schema.build-date=$BUILD_DATE \
    org.label-schema.name="gofsm" \
    org.label-schema.description="" \
    org.label-schema.url="" \
    org.label-schema.vcs-ref=$VCS_REF \
    org.label-schema.vcs-url="https://github.com/MrEhbr/gofsm" \
    org.label-schema.vendor="Alexey Burmistrov" \
    org.label-schema.version=$VERSION \
    org.label-schema.schema-version="1.0" \
    org.label-schema.cmd="docker run -i -t --rm MrEhbr/gofsm" \
    org.label-schema.help="docker exec -it $CONTAINER gofsm --help"
COPY            --from=builder /go/bin/gofsm /bin/
ENTRYPOINT      ["/bin/gofsm"]
#CMD             []
