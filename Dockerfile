FROM teamserverless/license-check:0.3.9 as license-check

FROM golang:1.15 as build

ARG GIT_COMMIT
ARG GIT_COMMIT_MESSAGE
ARG VERSION

COPY --from=license-check /license-check /usr/bin/

WORKDIR /go/src/github.com/nitishkumar71/gcp-cleaner

COPY pkg         pkg
COPY go.sum        go.sum
COPY go.mod        go.mod
COPY pkg           pkg
COPY main.go        .
RUN license-check -path ./ --verbose=false "Nitishkumar Singh"

# RUN go test -cover

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/cleaner .

FROM --platform=${TARGETPLATFORM:-linux/amd64} alpine:3.12 as ship

LABEL org.label-schema.license="MIT" \
    org.label-schema.vcs-url="https://github.com/nitishkumar71/gcp-cleaner" \
    org.label-schema.vcs-type="Git" \
    org.label-schema.name="nitishkumar71/gcp-cleaner" \
    org.label-schema.vendor="Nitishkumar Singh" \
    org.label-schema.docker.schema-version="1.0"

RUN addgroup -S app \
    && adduser -S -g app app \
    && apk add --no-cache ca-certificates

WORKDIR /home/app
EXPOSE 8080
ENV PORT 8080
ENV GIN_MODE release

COPY --from=build /go/src/github.com/nitishkumar71/gcp-cleaner/bin/cleaner .

RUN chown -R app:app ./
USER app
CMD ["./cleaner"]