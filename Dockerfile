FROM golang:1 as build

ARG VERSION="dev"

WORKDIR /TySug
COPY . .

RUN go test -test.short -test.v -test.race ./...
RUN CGO_ENABLED=0 GO111MODULE=on go build -v -a -ldflags "-w -X main.Version=${VERSION}" ./cmd/web

FROM gcr.io/distroless/base as base

FROM scratch

ARG VERSION="dev"
ARG GIT_REF="none"

LABEL org.label-schema.description="The TySug webservice Docker image. Suggesting typo-alternatives" \
      org.label-schema.schema-version="1.0" \
      org.label-schema.url="https://tysug.net/" \
      org.label-schema.vcs-url="https://github.com/Dynom/TySug" \
      org.label-schema.vcs-ref="${GIT_REF}" \
      org.label-schema.version="${VERSION}"

COPY --from=base ["/etc/ssl/certs/ca-certificates.crt", "/etc/ssl/certs/ca-certificates.crt"]
COPY --from=base ["/usr/share/zoneinfo", "/usr/share/zoneinfo"]
COPY --from=build ["/TySug/web", "/tysug"]
COPY --from=build ["/TySug/cmd/web/config.toml", "/"]

# Takes presedence over the configuration.
ENV LISTEN_URL="0.0.0.0:1337"
EXPOSE 1337


ENTRYPOINT ["/tysug"]
