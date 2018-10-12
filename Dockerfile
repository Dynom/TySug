FROM golang:1 as build

ARG VERSION

WORKDIR /go/src/github.com/Dynom/TySug
COPY . .

RUN go test -test.short -test.v -test.race ./...
RUN CGO_ENABLED=0 go install -v -a -ldflags "-w -X main.Version=${VERSION}" ./...

FROM gcr.io/distroless/base as base

FROM scratch

ARG VERSION
ARG GIT_REF

LABEL org.label-schema.description="The TySug webservice Docker image. Suggesting typo-alternatives" \
      org.label-schema.schema-version="1.0" \
      org.label-schema.url="https://tysug.net/" \
      org.label-schema.vcs-url="https://github.com/Dynom/TySug" \
      org.label-schema.vcs-ref="${GIT_REF}" \
      org.label-schema.version="${VERSION}"

COPY --from=base ["/etc/ssl/certs/ca-certificates.crt", "/etc/ssl/certs/ca-certificates.crt"]
COPY --from=base ["/usr/share/zoneinfo", "/usr/share/zoneinfo"]
COPY --from=build /go/bin/TySug /
COPY ["config.yml", "/"]

CMD ["/TySug"]